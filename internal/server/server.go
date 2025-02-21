package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/masintxi/gamehub/internal/auth"
	"github.com/masintxi/gamehub/internal/client"
	"github.com/masintxi/gamehub/internal/config"
	"github.com/masintxi/gamehub/internal/handlers"
)

type Server struct {
	router    chi.Router
	client    *client.Client
	steamAuth *auth.SteamAuth
	handlers  *handlers.SteamHandlers
	port      string
	domain    string
}

func NewServer(client *client.Client, steamAuth *auth.SteamAuth, cfg *config.Config) *Server {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	handlers := handlers.NewSteamHandlers(client, steamAuth)

	server := &Server{
		router:    r,
		client:    client,
		steamAuth: steamAuth,
		handlers:  handlers,
		port:      cfg.Server.Port,
		domain:    cfg.Server.Domain,
	}

	// Setup routes
	server.SetupRoutes()

	return server
}

func (s *Server) Start() error {
	fmt.Println("Starting server on :" + s.port)
	addr := fmt.Sprintf("%s:%s", s.domain, s.port)
	return http.ListenAndServe(addr, s.router)
}
