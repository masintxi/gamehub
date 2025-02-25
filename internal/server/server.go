package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/masintxi/gamehub/internal/auth"
	"github.com/masintxi/gamehub/internal/client"
	"github.com/masintxi/gamehub/internal/handlers"
)

type Server struct {
	Router    chi.Router
	Client    *client.Client
	SteamAuth *auth.SteamAuth
	Handlers  *handlers.SteamHandlers
	Port      string
	Domain    string
}

func NewServer(client *client.Client, steamAuth *auth.SteamAuth, server *Server) *Server {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	handlers := handlers.NewSteamHandlers(client, steamAuth)

	server.Router = r
	server.Client = client
	server.SteamAuth = steamAuth
	server.Handlers = handlers

	// Setup routes
	server.SetupRoutes()

	return server
}

func (s *Server) Start() error {
	fmt.Println("Starting server on :" + s.Port)
	addr := fmt.Sprintf("%s:%s", s.Domain, s.Port)
	return http.ListenAndServe(addr, s.Router)
}
