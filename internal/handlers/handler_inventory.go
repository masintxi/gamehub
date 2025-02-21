package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// const (
// 	providerName = "steam"
// 	authPath     = "/steam"
// 	callbackPath = "/steam/callback"
// 	sessionName  = "steam-session"
// )

func (h *SteamHandlers) HandleTradeInventory(w http.ResponseWriter, r *http.Request) {
	// 1. Get session and validate
	steamID, err := h.steamAuth.GetSteamID(r)
	if err != nil {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	// 2. Make the request
	url := fmt.Sprintf("https://api.steampowered.com/IEconService/GetInventoryItemsWithDescriptions/v1/?key=%s&steamid=%s&appid=753&contextid=6&get_descriptions=true",
		h.steamAuth.GetAPIKey(), steamID)

	headers := map[string]string{
		"User-Agent": "Mozilla/5.0",
	}

	bodyBytes, err := h.getResponseBody(url, headers)
	if err != nil {
		log.Printf("Error reading response: %v", err)
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	// 5. Decode JSON into struct
	var tradeResponse SteamTradeResponse
	if err := json.Unmarshal(bodyBytes, &tradeResponse); err != nil {
		log.Printf("Error decoding inventory: %v", err)
		http.Error(w, "Failed to decode inventory", http.StatusInternalServerError)
		return
	}

	// Set response headers and send response
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(bodyBytes); err != nil {
		log.Printf("Error writing response: %v", err)
		return
	}
}

func (h *SteamHandlers) HandleInventory(w http.ResponseWriter, r *http.Request) {
	// 1. Get session and validate
	steamID, err := h.steamAuth.GetSteamID(r)
	if err != nil {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	// 2. Make the request
	url := fmt.Sprintf("https://steamcommunity.com/inventory/%s/753/6?l=english&count=10", steamID)

	headers := map[string]string{
		"User-Agent": "Mozilla/5.0",
	}

	// 4. Read body
	bodyBytes, err := h.getResponseBody(url, headers)
	if err != nil {
		log.Printf("Error reading response: %v", err)
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	// 5. Decode JSON into struct
	var inventory SteamInventoryResponse
	if err := json.Unmarshal(bodyBytes, &inventory); err != nil {
		log.Printf("Error decoding inventory: %v", err)
		http.Error(w, "Failed to decode inventory", http.StatusInternalServerError)
		return
	}

	// 7. Send response
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(bodyBytes); err != nil {
		log.Printf("Error writing response: %v", err)
		return
	}
}
