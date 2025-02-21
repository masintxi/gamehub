package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
)

type OwnedGames struct {
	Response struct {
		GameCount int            `json:"game_count"`
		Games     []GameFromList `json:"games"`
	} `json:"response"`
}

func (h *SteamHandlers) HandleUserGames(w http.ResponseWriter, r *http.Request) {
	// 1. Get session and validate
	steamID, err := h.steamAuth.GetSteamID(r)
	if err != nil {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	// 2. Make the request
	url := fmt.Sprintf("https://api.steampowered.com/IPlayerService/GetOwnedGames/v1/?key=%s&steamid=%s&include_appinfo=true&include_extended_appinfo=true&format=json",
		h.steamAuth.GetAPIKey(), steamID)

	headers := map[string]string{
		"User-Agent": "Mozilla/5.0",
	}

	bodyBytes, err := h.getResponseBody(url, headers)
	if err != nil {
		log.Printf("Error fetching inventory: %v", err)
		http.Error(w, "Failed to fetch inventory", http.StatusInternalServerError)
		return
	}

	var ownedGamesList OwnedGames
	err = json.Unmarshal(bodyBytes, &ownedGamesList)
	if err != nil {
		fmt.Println("Error unmarshalling response:", err)
		return
	}

	games := make([]GameFromList, len(ownedGamesList.Response.Games))
	for i, game := range ownedGamesList.Response.Games {
		games[i] = GameFromList{
			AppID:                game.AppID,
			Name:                 game.Name,
			PlaytimeForever:      game.PlaytimeForever,
			PlaytimeDisconnected: game.PlaytimeDisconnected,
			Playtime2Weeks:       game.Playtime2Weeks,
		}
	}

	if len(games) == 0 {
		fmt.Println("No games found")
		return
	}

	fmt.Printf("%d games found\n", len(games))

	sort.SliceStable(games, func(i, j int) bool {
		return games[i].PlaytimeForever > games[j].PlaytimeForever
	})

	for i := 0; i < 10 && i < len(games); i++ {
		game := h.GetGameData(strconv.FormatInt(int64(games[i].AppID), 10))
		fmt.Printf("- %s (AppID: %d, Tiempo jugado: %d minutos)\n",
			games[i].Name, games[i].AppID, games[i].PlaytimeForever)
		fmt.Printf("  * %s\n", game.Data.ShortDescription)
	}
}
