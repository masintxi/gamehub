package auth

import (
	"fmt"
	"log"
	"net/http"

	"github.com/markbates/goth"
)

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
	// 1. Get or create session
	cookieSession, err := sa.Store.Get(r, sessionName)
	if err != nil {
		log.Printf("Error getting session: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 2. Handle Steam authentication
	user, err := sa.authenticateWithSteam(r)
	if err != nil {
		log.Printf("Steam authentication error: %v", err)
		http.Error(w, "Authentication Failed", http.StatusInternalServerError)
		return
	}

	// 3. Store user info in session
	cookieSession.Values["steamID"] = user.UserID
	cookieSession.Values["userName"] = user.Name

	// Initialize market session
	err = sa.InitializeMarketSession()
	if err != nil {
		log.Printf("Failed to initialize market session: %v", err)
	}

	// 4. Save session
	if err := cookieSession.Save(r, w); err != nil {
		log.Printf("Error saving session: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 5. Send response
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{
        "steamID": "%s",
        "name": "%s",
        "sessionCaptured": false
    }`, user.UserID, user.Name)
}

func (sa *SteamAuth) authenticateWithSteam(r *http.Request) (goth.User, error) {
	gothProvider, err := goth.GetProvider(providerName)
	if err != nil {
		return goth.User{}, fmt.Errorf("error getting provider: %w", err)
	}

	gothSession, err := gothProvider.BeginAuth(r.URL.Query().Get("state"))
	if err != nil {
		return goth.User{}, fmt.Errorf("error beginning auth: %w", err)
	}

	_, err = gothSession.Authorize(gothProvider, r.URL.Query())
	if err != nil {
		return goth.User{}, fmt.Errorf("error authorizing: %w", err)
	}

	user, err := gothProvider.FetchUser(gothSession)
	if err != nil {
		return goth.User{}, fmt.Errorf("error fetching user: %w", err)
	}

	return user, nil
}
