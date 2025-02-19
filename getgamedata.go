package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (cfg *apiConfig) GetGameNameByID(gameID string) string {
	return cfg.GetGameData(gameID).Data.Name
}

func (cfg *apiConfig) GetGameData(gameID string) GameData {
	url := fmt.Sprintf("https://store.steampowered.com/api/appdetails?appids=%s", gameID)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return GameData{}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return GameData{}
	}

	var gameData map[string]GameData
	err = json.Unmarshal(body, &gameData)
	if err != nil {
		fmt.Println("Error unmarshalling response:", err)
		return GameData{}
	}

	//fmt.Println("Response:", string(body))
	//fmt.Println(gameData)
	//fmt.Println("Game name:", gameData[gameID].Data.Name)

	return gameData[gameID]
}

func (cfg *apiConfig) GetGameStats(gameID string) {
	//url := fmt.Sprintf("https://store.steampowered.com/api/appdetails?appids=%s&l=english&cc=US&filters=priceoverview", gameID)
	//url := fmt.Sprintf("https://api.steampowered.com/ISteamUserStats/GetNumberOfCurrentPlayers/v1/?appid=%s", gameID)
	url := fmt.Sprintf("https://api.steampowered.com/ISteamUserStats/GetSchemaForGame/v2/?key=%s&appid=%s", cfg.apiKey, gameID)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	fmt.Println("Response:", string(body))
}
