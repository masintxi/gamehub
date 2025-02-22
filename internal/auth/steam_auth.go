package auth

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/steam"
)

type SteamAuth struct {
	ApiKey           string
	CallbackURL      string
	Provider         *steam.Provider
	Store            *sessions.CookieStore
	SessionID        string
	SteamLoginSecure string
}

const (
	providerName = "steam"
	sessionName  = "steam-session"
)

func NewSteamAuth(apiKey, callbackURL string) *SteamAuth {
	// Create a cookie store with a secret key
	store := sessions.NewCookieStore([]byte("your-secret-key"))

	// Initialize the Steam provider
	provider := steam.New(apiKey, callbackURL)
	goth.UseProviders(provider)

	return &SteamAuth{
		ApiKey:      apiKey,
		CallbackURL: callbackURL,
		Provider:    provider,
		Store:       store,
	}
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
