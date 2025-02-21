package handlers

import (
	"encoding/json"
	"fmt"
	"log"
)

func (h *SteamHandlers) GetGameData(gameID string) GameData {

	// 2. Make the request
	url := fmt.Sprintf("https://store.steampowered.com/api/appdetails?appids=%s", gameID)

	headers := map[string]string{
		"User-Agent": "Mozilla/5.0",
	}

	bodyBytes, err := h.getResponseBody(url, headers)
	if err != nil {
		log.Printf("Error fetching inventory: %v", err)
		return GameData{}
	}

	var gameData map[string]GameData
	err = json.Unmarshal(bodyBytes, &gameData)
	if err != nil {
		log.Println("Error unmarshalling response:", err)
		return GameData{}
	}

	return gameData[gameID]
}

func (h *SteamHandlers) GetGameStats(gameID string) {
	//url := fmt.Sprintf("https://store.steampowered.com/api/appdetails?appids=%s&l=english&cc=US&filters=priceoverview", gameID)
	//url := fmt.Sprintf("https://api.steampowered.com/ISteamUserStats/GetNumberOfCurrentPlayers/v1/?appid=%s", gameID)
	url := fmt.Sprintf("https://api.steampowered.com/ISteamUserStats/GetSchemaForGame/v2/?key=%s&appid=%s", h.steamAuth.GetAPIKey(), gameID)

	headers := map[string]string{
		"User-Agent": "Mozilla/5.0",
	}

	bodyBytes, err := h.getResponseBody(url, headers)
	if err != nil {
		log.Printf("Error fetching inventory: %v", err)
		return
	}

	fmt.Println("Response:", string(bodyBytes))
}
