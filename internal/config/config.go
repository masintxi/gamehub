package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// Steam Configuration
	SteamAPIKey      string
	SteamUserID      string
	SteamCallbackURL string

	// Cache Configuration
	Cache struct {
		ProjectName     string
		CleanupInterval time.Duration
		MaxSize         int64
		FileExtension   string
		ExpireAfter     time.Duration
		Compression     bool
		CachePath       string
	}

	// Server Configuration
	Server struct {
		Port   string
		Domain string
	}
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	cfg := &Config{}

	// Load Steam config
	cfg.SteamAPIKey = os.Getenv("STEAM_API_KEY")
	if cfg.SteamAPIKey == "" {
		log.Fatal("STEAM_API_KEY not set")
	}
	cfg.SteamUserID = os.Getenv("STEAM_USER_ID")
	if cfg.SteamUserID == "" {
		log.Fatal("STEAM_USER_ID not set")
	}
	cfg.SteamCallbackURL = "http://localhost:8080/auth/steam/callback"

	// Set cache config
	cfg.Cache.ProjectName = "gamehub"
	cfg.Cache.CleanupInterval = 5 * time.Second
	cfg.Cache.MaxSize = 1024 * 1024 * 10 // 10 MB
	cfg.Cache.FileExtension = "json"
	cfg.Cache.ExpireAfter = 30 * time.Minute
	cfg.Cache.Compression = true

	// Set server config
	cfg.Server.Port = os.Getenv("SERVER_PORT")
	if cfg.Server.Port == "" {
		cfg.Server.Port = "8080"
	}
	cfg.Server.Domain = os.Getenv("SERVER_DOMAIN")
	if cfg.Server.Domain == "" {
		cfg.Server.Domain = "localhost"
	}

	return cfg
}
