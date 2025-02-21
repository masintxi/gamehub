package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (h *SteamHandlers) HandleUserData(w http.ResponseWriter, r *http.Request) {
	// 1. Get session and validate
	steamID, err := h.steamAuth.GetSteamID(r)
	if err != nil {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}
	// 2. Make the request
	url := fmt.Sprintf("https://api.steampowered.com/ISteamUser/GetPlayerSummaries/v0002/?key=%s&steamids=%s",
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

	// 3. Decode JSON into struct
	var playerResponse PlayerResponse
	if err := json.Unmarshal(bodyBytes, &playerResponse); err != nil {
		log.Printf("Error decoding inventory: %v", err)
		http.Error(w, "Failed to decode inventory", http.StatusInternalServerError)
		return
	}
	// 5. Send response
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(bodyBytes); err != nil {
		log.Printf("Error writing response: %v", err)
		return
	}

	fmt.Println("Response:", string(bodyBytes))

}
