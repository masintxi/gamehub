package main

import (
	"log"

	"github.com/masintxi/gamehub/internal/auth"
	"github.com/masintxi/gamehub/internal/cache"
	"github.com/masintxi/gamehub/internal/client"
	"github.com/masintxi/gamehub/internal/config"
	"github.com/masintxi/gamehub/internal/server"
)

func main() {
	cfg := config.Load()

	client := client.NewClient(cache.CacheConfig{
		ProjectName:     cfg.Cache.ProjectName,
		CleanupInterval: cfg.Cache.CleanupInterval,
		MaxSize:         cfg.Cache.MaxSize,
		FileExtension:   cfg.Cache.FileExtension,
		ExpireAfter:     cfg.Cache.ExpireAfter,
		Compression:     cfg.Cache.Compression,
		CachePath:       cfg.Cache.CachePath,
	})

	steamAuth := auth.NewSteamAuth(
		cfg.SteamAPIKey,
		cfg.SteamCallbackURL,
	)

	server := server.NewServer(client, steamAuth, cfg)

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
