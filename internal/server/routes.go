package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

const (
	// providerName = "steam"
	authPath     = "/steam"
	callbackPath = "/steam/callback"
	// sessionName  = "steam-session"
)

func (s *Server) SetupRoutes() {
	s.router.Get("/", s.HandleHome)
	s.router.Route("/auth", func(r chi.Router) {
		r.Get(authPath, s.steamAuth.HandleSteamLogin)
		r.Get(callbackPath, s.steamAuth.HandleSteamCallback)
	})
	s.router.Get("/inventory", s.handlers.HandleInventory)
	s.router.Get("/trade-inventory", s.handlers.HandleTradeInventory)
	s.router.Get("/market/{market_hash_name}", s.handlers.HandleMarketData)
	s.router.Get("/user-data", s.handlers.HandleUserData)
	s.router.Get("/user-games", s.handlers.HandleUserGames)
}

func (s *Server) HandleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `
        <html>
            <body>
                <a href="/auth/steam">Login with Steam</a>
                <br/>
                <a href="/trade-inventory">View Inventory (Trading API)</a>
				<br/>
				<a href="/user-data">View User Data</a>
				<br/>
				<a href="/user-games">View User Games</a>
				<br/>
            </body>
        </html>
    `)
}
