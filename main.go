package main

import (
	"time"

	"github.com/masintxi/gamehub/internal/api"
	"github.com/masintxi/gamehub/internal/cache"
)

type Config struct {
	client *api.Client
}

func main() {

	cfg := Config{
		client: api.NewClient(cache.CacheConfig{
			ProjectName:     "gamehub",
			CleanupInterval: 5 * time.Second,
			MaxSize:         1024 * 1024 * 10, // 10 MB
			FileExtension:   "json",
			ExpireAfter:     30 * time.Minute,
			Compression:     true,
			CachePath:       "",
		}),
	}

	cfg.client.GetUserGames()

}

//cfg.GetUserData()
//cfg.GetGameData("2457220")
//cfg.GetGameStats("2457220")
