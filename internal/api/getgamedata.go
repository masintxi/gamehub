package api

import (
	"encoding/json"
	"fmt"
	"log"
)

func (client *Client) GetGameData(gameID string) GameData {
	url := fmt.Sprintf("https://store.steampowered.com/api/appdetails?appids=%s", gameID)

	body, err := client.getResposeBody(url)
	if err != nil {
		log.Println(err)
		return GameData{}
	}

	var gameData map[string]GameData
	err = json.Unmarshal(body, &gameData)
	if err != nil {
		log.Println("Error unmarshalling response:", err)
		return GameData{}
	}

	return gameData[gameID]
}

func (client *Client) GetGameStats(gameID string) {
	//url := fmt.Sprintf("https://store.steampowered.com/api/appdetails?appids=%s&l=english&cc=US&filters=priceoverview", gameID)
	//url := fmt.Sprintf("https://api.steampowered.com/ISteamUserStats/GetNumberOfCurrentPlayers/v1/?appid=%s", gameID)
	url := fmt.Sprintf("https://api.steampowered.com/ISteamUserStats/GetSchemaForGame/v2/?key=%s&appid=%s", client.apiKey, gameID)

	body, err := client.getResposeBody(url)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("Response:", string(body))
}
