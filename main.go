package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type apiConfig struct {
	apiKey  string
	steamID string
}

func main() {
	// Replace with your Steam API key and SteamID
	godotenv.Load()
	apiKey := os.Getenv("STEAM_API_KEY")
	if apiKey == "" {
		log.Fatal("STEAM_API_KEY not setted. Check .env file")
	}
	steamID := os.Getenv("STEAM_USER_ID")
	if steamID == "" {
		log.Fatal("STEAM_USER_ID not setted. Check .env file")
	}

	cfg := apiConfig{
		apiKey:  apiKey,
		steamID: steamID,
	}

	//cfg.GetUserData()
	//cfg.GetGameData("2457220")
	//cfg.GetGameStats("2457220")
	cfg.GetUserGames()
}
