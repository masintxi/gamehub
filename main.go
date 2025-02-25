package main

import (
	"log"

	"github.com/masintxi/gamehub/internal/auth"
	"github.com/masintxi/gamehub/internal/client"
	"github.com/masintxi/gamehub/internal/config"
	"github.com/masintxi/gamehub/internal/server"
)

func main() {
	cfg := config.Load()

	client := client.NewClient(cfg.CacheConfig)

	steamAuth := auth.NewSteamAuth(cfg.SteamAuth)

	server := server.NewServer(client, steamAuth, &cfg.Server)

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}

// fmt.Println("Starting server on :8080")
// log.Fatal(http.ListenAndServe(":8080", r))

//cfg.client.GetUserGames()
//cfg.client.ListInventory()
//cfg.GetUserData()
//cfg.GetGameData("2457220")
//cfg.GetGameStats("2457220")
