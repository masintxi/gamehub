package auth

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/steam"
)

type SteamAuth struct {
	apiKey      string
	callbackURL string
	provider    *steam.Provider
	store       *sessions.CookieStore
}

const (
	providerName = "steam"
	// authPath     = "/steam"
	// callbackPath = "/steam/callback"
	sessionName = "steam-session"
)

func NewSteamAuth(apiKey, callbackURL string) *SteamAuth {
	// Create a cookie store with a secret key
	store := sessions.NewCookieStore([]byte("your-secret-key"))

	// Initialize the Steam provider
	provider := steam.New(apiKey, callbackURL)
	goth.UseProviders(provider)

	return &SteamAuth{
		apiKey:      apiKey,
		callbackURL: callbackURL,
		provider:    provider,
		store:       store,
	}
}

func (sa *SteamAuth) GetSteamID(r *http.Request) (string, error) {
	session, err := sa.store.Get(r, sessionName)
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
	return sa.store.Get(r, sessionName)
}

func (sa *SteamAuth) GetAPIKey() string {
	return sa.apiKey
}

func (sa *SteamAuth) HandleSteamLogin(w http.ResponseWriter, r *http.Request) {
	// Get the Steam provider
	gothProvider, err := goth.GetProvider(providerName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Begin the authentication session
	session, err := gothProvider.BeginAuth(r.URL.Query().Get("state"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the authentication URL
	url, err := session.GetAuthURL()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to Steam's login page
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (sa *SteamAuth) HandleSteamCallback(w http.ResponseWriter, r *http.Request) {
	cookieSession, err := sa.store.Get(r, sessionName)
	if err != nil {
		log.Printf("Error getting session: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Get the Steam provider
	gothProvider, err := goth.GetProvider(providerName)
	if err != nil {
		log.Printf("Error getting provider: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Begin the auth session with Steam
	gothSession, err := gothProvider.BeginAuth(r.URL.Query().Get("state"))
	if err != nil {
		log.Printf("Error beginning auth: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Complete the authentication process
	_, err = gothSession.Authorize(gothProvider, r.URL.Query())
	if err != nil {
		log.Printf("Error authorizing: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Fetch user info
	user, err := gothProvider.FetchUser(gothSession)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Store user info in cookie session
	cookieSession.Values["steamID"] = user.UserID
	cookieSession.Values["userName"] = user.Name

	err = cookieSession.Save(r, w)
	if err != nil {
		log.Printf("Error saving session: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// At this point, user.UserID will contain the Steam ID
	// You might want to store this in a session or return it to the client
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"steamID": "%s", "name": "%s"}`, user.UserID, user.Name)
}

// func (sa *SteamAuth) HandleHome(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, `
//         <html>
//             <body>
//                 <a href="/auth/steam">Login with Steam</a>
//                 <br/>
//                 <a href="/trade-inventory">View Inventory (Trading API)</a>
//             </body>
//         </html>
//     `)
// }
