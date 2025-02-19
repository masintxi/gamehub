package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
)

type Game struct {
	AppID           int    `json:"appid"`
	Name            string `json:"name"`
	PlaytimeForever int    `json:"playtime_forever"`
}

type OwnedGames struct {
	Response struct {
		GameCount int    `json:"game_count"`
		Games     []Game `json:"games"`
	} `json:"response"`
}

func (cfg *apiConfig) GetUserGames() {
	url := fmt.Sprintf("https://api.steampowered.com/IPlayerService/GetOwnedGames/v1/?key=%s&steamid=%s", cfg.apiKey, cfg.steamID)
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

	var ownedGamesList OwnedGames
	err = json.Unmarshal(body, &ownedGamesList)
	if err != nil {
		fmt.Println("Error unmarshalling response:", err)
		return
	}

	games := make([]Game, len(ownedGamesList.Response.Games))
	for i, game := range ownedGamesList.Response.Games {
		games[i] = Game{
			AppID:           game.AppID,
			Name:            game.Name,
			PlaytimeForever: game.PlaytimeForever,
		}
	}

	sort.SliceStable(games, func(i, j int) bool {
		return games[i].PlaytimeForever > games[j].PlaytimeForever
	})

	for i := 0; i < 5 && i < len(games); i++ {
		game := cfg.GetGameData(strconv.FormatInt(int64(games[i].AppID), 10))
		fmt.Printf("- %s (AppID: %d, Tiempo jugado: %d minutos)\n", game.Data.Name, games[i].AppID, games[i].PlaytimeForever)
		fmt.Printf("  * %s\n", game.Data.Genres)
	}

	//fmt.Println("Response:", string(body))
}
