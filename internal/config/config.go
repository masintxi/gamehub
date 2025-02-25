package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/masintxi/gamehub/internal/auth"
	"github.com/masintxi/gamehub/internal/cache"
	"github.com/masintxi/gamehub/internal/server"
)

type Config struct {
	// Server Configuration
	Server server.Server

	// Steam Configuration
	SteamAuth auth.SteamAuth

	// Cache Configuration
	CacheConfig cache.CacheConfig
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	cfg := &Config{}

	// Set server config
	cfg.Server.Port = os.Getenv("SERVER_PORT")
	if cfg.Server.Port == "" {
		cfg.Server.Port = "8080"
	}
	cfg.Server.Domain = os.Getenv("SERVER_DOMAIN")
	if cfg.Server.Domain == "" {
		cfg.Server.Domain = "localhost"
	}

	// Load Steam config
	cfg.SteamAuth.ApiKey = os.Getenv("STEAM_API_KEY")
	if cfg.SteamAuth.ApiKey == "" {
		log.Fatal("STEAM_API_KEY not set")
	}
	// cfg.SteamAuth.SteamUserID = os.Getenv("STEAM_USER_ID")
	// if cfg.SteamAuth.SteamUserID == "" {
	// 	log.Fatal("STEAM_USER_ID not set")
	// }

	cfg.SteamAuth.CallbackURL = fmt.Sprintf("http://%s:%s/auth/steam/callback",
		cfg.Server.Domain, cfg.Server.Port)

	// Set cache config
	cfg.CacheConfig.ProjectName = "gamehub"
	cfg.CacheConfig.CleanupInterval = 5 * time.Second
	cfg.CacheConfig.MaxSize = 1024 * 1024 * 10 // 10 MB
	cfg.CacheConfig.FileExtension = "json"
	cfg.CacheConfig.Compression = true
	cfg.CacheConfig.ExpireAfter = 30 * time.Minute

	return cfg
}
