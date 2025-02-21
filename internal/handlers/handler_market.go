package handlers

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

type MarketPriceHistory struct {
	Success      bool            `json:"success"`
	PriceHistory [][]interface{} `json:"price_history"`
	// Each price history entry is [timestamp, price, volume]
}

func (h *SteamHandlers) HandleMarketData(w http.ResponseWriter, r *http.Request) {
	// First, get the user's session
	session, err := h.steamAuth.GetSession(r)
	if err != nil {
		log.Printf("Session error: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if user is authenticated
	if session.Values["steamID"] == nil {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}

	marketHashName := chi.URLParam(r, "market_hash_name")
	if marketHashName == "" {
		http.Error(w, "Market hash name is required", http.StatusBadRequest)
		return
	}
	marketHashName = url.QueryEscape(marketHashName)
	marketHashName = strings.ReplaceAll(marketHashName, "+", "%20")

	// Steam Market API URL for price history
	// url := fmt.Sprintf("https://steamcommunity.com/market/pricehistory/?appid=753&market_hash_name=%s", url.QueryEscape(marketHashName))

	//baseURL := fmt.Sprintf("https://steamcommunity.com/market/listings/753/%s/render/history", marketHashName)
	baseURL := "https://steamcommunity.com/market/pricehistory/"
	u, err := url.Parse(baseURL)
	if err != nil {
		http.Error(w, "Failed to parse URL", http.StatusInternalServerError)
		return
	}

	q := u.Query()
	q.Add("appid", "753")
	q.Add("market_hash_name", marketHashName)
	u.RawQuery = q.Encode()

	// Add more detailed debug logging
	log.Printf("Final URL: %s", u.String())
	log.Printf("Market Hash Name (raw): %s", chi.URLParam(r, "market_hash_name"))
	log.Printf("Market Hash Name (encoded): %s", marketHashName)

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	for _, cookie := range r.Cookies() {
		req.AddCookie(cookie)
	}
	// Add more complete headers
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	// Add the specific Steam cookies
	if sessionID, ok := session.Values["sessionid"].(string); ok {
		req.AddCookie(&http.Cookie{
			Name:  "sessionid",
			Value: sessionID,
		})
	}

	// Add steamCountry cookie
	if steamCountry, ok := session.Values["steamCountry"].(string); ok {
		req.AddCookie(&http.Cookie{
			Name:  "steamCountry",
			Value: steamCountry,
		})
	}

	// Add timezoneOffset cookie
	req.AddCookie(&http.Cookie{
		Name:  "timezoneOffset",
		Value: "3600,0",
	})

	// Create a client that handles gzip compression
	client := &http.Client{
		Transport: &http.Transport{
			DisableCompression: false,
		},
		Timeout: 30 * time.Second,
	}

	log.Printf("Requesting URL: %s", u.String())
	log.Printf("With cookies: %v", req.Cookies())

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to fetch market data", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Steam Error - Status: %d", resp.StatusCode)
		log.Printf("Steam Error - Headers: %v", resp.Header)
		log.Printf("Steam Error - URL: %s", baseURL)
	}

	var body []byte
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err := gzip.NewReader(resp.Body)
		if err != nil {
			log.Printf("Error creating gzip reader: %v", err)
			return
		}
		defer reader.Close()
		body, err = io.ReadAll(reader)
		if err != nil {
			log.Printf("Error reading gzipped body: %v", err)
			return
		}
	} else {
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			return
		}
	}

	log.Printf("Response Status: %d", resp.StatusCode)
	log.Printf("Response Body: %s", string(body))

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
