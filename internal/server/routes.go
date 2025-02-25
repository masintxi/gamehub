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
	s.Router.Get("/", s.HandleHome)
	s.Router.Route("/auth", func(r chi.Router) {
		r.Get(authPath, s.SteamAuth.HandleSteamLogin)
		r.Get(callbackPath, s.SteamAuth.HandleSteamCallback)
	})
	s.Router.Get("/inventory", s.Handlers.HandleInventory)
	s.Router.Get("/trade-inventory", s.Handlers.HandleTradeInventory)
	s.Router.Get("/market/{market_hash_name}", s.Handlers.HandleMarketData)
	s.Router.Get("/user-data", s.Handlers.HandleUserData)
	s.Router.Get("/user-games", s.Handlers.HandleUserGames)
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
				<a href="/market/311690-Frifle and Mauser">Market Price for Frifle and Mauser</a>
				<br/>
            </body>
        </html>
    `)
}
