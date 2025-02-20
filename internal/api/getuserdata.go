package api

import (
	"fmt"
)

func (client *Client) GetUserData() {
	url := fmt.Sprintf("https://api.steampowered.com/ISteamUser/GetPlayerSummaries/v0002/?key=%s&steamids=%s",
		client.apiKey, client.steamID)

	body, err := client.getResposeBody(url)
	if err != nil {
		return
	}

	fmt.Println("Response:", string(body))
}
