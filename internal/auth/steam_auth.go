package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"

	"github.com/gorilla/sessions"
)

type SteamAuth struct {
	ApiKey           string
	CallbackURL      string
	Store            *sessions.CookieStore
	SessionID        string
	SteamLoginSecure string
	SteamClient      *http.Client
	//Provider         *steam.Provider
}

const (
	providerName = "steam"
	sessionName  = "steam-session"

	// Steam API Endpoints
	apiLoginEndpoint       = "https://steamcommunity.com/openid/login"
	apiUserSummaryEndpoint = "http://api.steampowered.com/ISteamUser/GetPlayerSummaries/v0002/?key=%s&steamids=%s"

	// OpenID settings
	openIDMode       = "checkid_setup"
	openIDNs         = "http://specs.openid.net/auth/2.0"
	openIDIdentifier = "http://specs.openid.net/auth/2.0/identifier_select"
)

func NewSteamAuth(sa SteamAuth) *SteamAuth {
	// Create a cookie store with a secret key
	store := sessions.NewCookieStore([]byte("your-secret-key"))
	sa.Store = store

	// Create a cookie jar
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("Failed to create cookie jar: %v", err)
	}

	// Attach the cookie jar to an HTTP client
	sa.SteamClient = &http.Client{
		Jar: jar,
	}

	return &sa
}

func (sa *SteamAuth) GetAuthURL() (string, error) {
	callbackURL, err := url.Parse(sa.CallbackURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse callback URL: %w", err)
	}

	urlValues := map[string]string{
		"openid.claimed_id": openIDIdentifier,
		"openid.identity":   openIDIdentifier,
		"openid.mode":       openIDMode,
		"openid.ns":         openIDNs,
		"openid.realm":      fmt.Sprintf("%s://%s", callbackURL.Scheme, callbackURL.Host),
		"openid.return_to":  callbackURL.String(),
	}

	u, err := url.Parse(apiLoginEndpoint)
	if err != nil {
		return "", fmt.Errorf("failed to parse login endpoint: %w", err)
	}

	v := u.Query()
	for key, value := range urlValues {
		v.Add(key, value)
	}
	u.RawQuery = v.Encode()

	return u.String(), nil
}

func (sa *SteamAuth) HandleLogin(w http.ResponseWriter, r *http.Request) {
	url, err := sa.GetAuthURL()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to Steam's login page
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (sa *SteamAuth) HandleCallback(w http.ResponseWriter, r *http.Request) {
	cookieSession, err := sa.Store.Get(r, sessionName)
	if err != nil {
		log.Printf("Error getting session: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	params := r.URL.Query()

	validationParams := url.Values{}
	validationParams.Add("openid.assoc_handle", params.Get("openid.assoc_handle"))
	validationParams.Add("openid.signed", params.Get("openid.signed"))
	validationParams.Add("openid.sig", params.Get("openid.sig"))
	validationParams.Add("openid.ns", params.Get("openid.ns"))
	validationParams.Add("openid.mode", "check_authentication")
	for _, param := range strings.Split(params.Get("openid.signed"), ",") {
		validationParams.Add("openid."+param, params.Get("openid."+param))
	}

	for key, values := range r.URL.Query() {
		log.Printf("Query Param: %s=%s\n", key, values)
	}
	resp, err := sa.SteamClient.PostForm(apiLoginEndpoint, validationParams)
	if err != nil {
		http.Error(w, "Error validating OpenID response", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading response body", http.StatusInternalServerError)
		return
	}

	response := strings.Split(string(body), "\n")
	if response[0] != "ns:"+openIDNs {
		http.Error(w, "Invalid OpenID response", http.StatusInternalServerError)
		return
	}

	if response[1] == "is_valid:false" {
		http.Error(w, "OpenID validation failed", http.StatusUnauthorized)
		return
	}

	openIDURL := params.Get("openid.claimed_id")
	validationRegExp := regexp.MustCompile("^(http|https)://steamcommunity.com/openid/id/[0-9]{15,25}$")
	if !validationRegExp.MatchString(openIDURL) {
		http.Error(w, "Invalid Steam ID pattern", http.StatusInternalServerError)
		return
	}
	log.Println("openIDURL: ", openIDURL)

	steamID := regexp.MustCompile(`\D+`).ReplaceAllString(openIDURL, "")
	ResponseNonce := params.Get("openid.response_nonce")
	log.Println("steamID: ", steamID)
	log.Println("ResponseNonce: ", ResponseNonce)

	steamURL, _ := url.Parse("https://steamcommunity.com")
	cookies := sa.SteamClient.Jar.Cookies(steamURL)
	for _, cookie := range cookies {
		log.Printf("Cookie in jar: %s=%s\n", cookie.Name, cookie.Value)
		//cookieSession.Values[cookie.Name] = cookie.Value
	}

	// for _, cookie := range resp.Cookies() {
	// 	log.Println("Response Cookie: ", cookie)
	// 	log.Printf("Cookie: Name=%s, Value=%s, Domain=%s, Path=%s, HttpOnly=%t, Secure=%t",
	// 		cookie.Name, cookie.Value, cookie.Domain, cookie.Path, cookie.HttpOnly, cookie.Secure)
	// }

	cookieSession.Values["steamID"] = steamID
	log.Println("cookieSession: ", cookieSession)

	steamUser, err := sa.FetchUser(cookieSession)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	cookieSession.Values["name"] = steamUser

	if err := cookieSession.Save(r, w); err != nil {
		log.Printf("Error saving session: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{
        "steamID": "%s",
        "name": "%s",
        "sessionCaptured": true
    }`, steamID, steamUser)
}

func (sa *SteamAuth) FetchUser(session *sessions.Session) (string, error) {
	apiResponse := struct {
		Response struct {
			Players []struct {
				Steamid                  string `json:"steamid"`
				Communityvisibilitystate int    `json:"communityvisibilitystate"`
				Profilestate             int    `json:"profilestate"`
				Personaname              string `json:"personaname"`
				Commentpermission        int    `json:"commentpermission"`
				Profileurl               string `json:"profileurl"`
				Avatar                   string `json:"avatar"`
				Avatarmedium             string `json:"avatarmedium"`
				Avatarfull               string `json:"avatarfull"`
				Avatarhash               string `json:"avatarhash"`
				Lastlogoff               int    `json:"lastlogoff"`
				Personastate             int    `json:"personastate"`
				Realname                 string `json:"realname"`
				Primaryclanid            string `json:"primaryclanid"`
				Timecreated              int    `json:"timecreated"`
				Personastateflags        int    `json:"personastateflags"`
			} `json:"players"`
		} `json:"response"`
	}{}

	steamID, ok := session.Values["steamID"].(string)
	if !ok {
		return "", fmt.Errorf("not authenticated")
	}

	userSummaryURL := fmt.Sprintf(apiUserSummaryEndpoint, sa.ApiKey, steamID)
	req, err := http.NewRequest("GET", userSummaryURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	fmt.Println(string(body))
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return "", err
	}

	if l := len(apiResponse.Response.Players); l != 1 {
		return "", fmt.Errorf("expected one player in API response, got %d", l)
	}

	return apiResponse.Response.Players[0].Personaname, nil
}

func (sa *SteamAuth) GetSteamID(r *http.Request) (string, error) {
	session, err := sa.Store.Get(r, sessionName)
	if err != nil {
		return "", err
	}

	steamID, ok := session.Values["steamID"].(string)
	if !ok {
		return "", fmt.Errorf("not authenticated")
	}

	return steamID, nil
}

func (sa *SteamAuth) GetSession(r *http.Request) (*sessions.Session, error) {
	return sa.Store.Get(r, sessionName)
}

func (sa *SteamAuth) GetAPIKey() string {
	return sa.ApiKey
}

func (sa *SteamAuth) SetSteamCookies(sessionID, steamLoginSecure string) {
	sa.SessionID = sessionID
	sa.SteamLoginSecure = steamLoginSecure
}

func (sa *SteamAuth) GetMarketCookies() map[string]string {
	return map[string]string{
		"sessionid":        sa.SessionID,
		"steamLoginSecure": sa.SteamLoginSecure,
	}
}
