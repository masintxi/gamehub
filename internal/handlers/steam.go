package handlers

import (
	"github.com/masintxi/gamehub/internal/auth"
	"github.com/masintxi/gamehub/internal/client"
)

type SteamHandlers struct {
	client    *client.Client
	steamAuth *auth.SteamAuth
}

func NewSteamHandlers(client *client.Client, steamAuth *auth.SteamAuth) *SteamHandlers {
	return &SteamHandlers{
		client:    client,
		steamAuth: steamAuth,
	}
}
