package api

// import (
// 	"encoding/json"
// 	"fmt"
// 	"sort"
// 	"strconv"
// )

// type OwnedGames struct {
// 	Response struct {
// 		GameCount int            `json:"game_count"`
// 		Games     []GameFromList `json:"games"`
// 	} `json:"response"`
// }

// func (client *Client) GetUserGames() {
// 	rootUrl := "https://api.steampowered.com/IPlayerService/GetOwnedGames/v1"
// 	url := fmt.Sprintf("%s/?key=%s&steamid=%s&include_appinfo=true&include_extended_appinfo=true&format=json",
// 		rootUrl, client.apiKey, client.steamID)

// 	body, err := client.getResposeBody(url)
// 	if err != nil {
// 		return
// 	}

// 	var ownedGamesList OwnedGames
// 	err = json.Unmarshal(body, &ownedGamesList)
// 	if err != nil {
// 		fmt.Println("Error unmarshalling response:", err)
// 		return
// 	}

// 	games := make([]GameFromList, len(ownedGamesList.Response.Games))
// 	for i, game := range ownedGamesList.Response.Games {
// 		games[i] = GameFromList{
// 			AppID:                game.AppID,
// 			Name:                 game.Name,
// 			PlaytimeForever:      game.PlaytimeForever,
// 			PlaytimeDisconnected: game.PlaytimeDisconnected,
// 			Playtime2Weeks:       game.Playtime2Weeks,
// 		}
// 	}

// 	if len(games) == 0 {
// 		fmt.Println("No games found")
// 		return
// 	}

// 	fmt.Printf("%d games found\n", len(games))

// 	sort.SliceStable(games, func(i, j int) bool {
// 		return games[i].PlaytimeForever > games[j].PlaytimeForever
// 	})

// 	for i := 0; i < 10 && i < len(games); i++ {
// 		game := client.GetGameData(strconv.FormatInt(int64(games[i].AppID), 10))
// 		fmt.Printf("- %s (AppID: %d, Tiempo jugado: %d minutos)\n",
// 			games[i].Name, games[i].AppID, games[i].PlaytimeForever)
// 		fmt.Printf("  * %s\n", game.Data.ShortDescription)
// 	}

// 	//fmt.Println("Response:", string(body))
// }
