package handlers

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/go-chi/chi"
)

type MarketPriceHistory struct {
	Success      bool            `json:"success"`
	PriceHistory [][]interface{} `json:"price_history"`
	// Each price history entry is [timestamp, price, volume]
}

func (h *SteamHandlers) HandleMarketItem(w http.ResponseWriter, r *http.Request) {
	itemName := chi.URLParam(r, "item_name")
	log.Println("itemName: ", itemName)
	if itemName == "" {
		itemName = "Frifle and Mauser"
	}

	url := fmt.Sprintf("https://steamcommunity.com/market/listings/753/%s", url.QueryEscape(itemName))
	log.Println(url)

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(`/mnt/c/Program Files/Google/Chrome/Application/chrome.exe`),
		chromedp.Headless,                     // Run in headless mode
		chromedp.Flag("enable-logging", "v1"), // Enable logging
	)

	// Create a context with the custom options
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// Create a new browser context
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Set a timeout for the operation
	ctx, cancel = context.WithTimeout(ctx, 60*time.Second) // Increase timeout
	defer cancel()

	// Variable to store the extracted price history
	var priceHistory string

	// Run the chromedp tasks
	err := chromedp.Run(ctx,
		// Navigate to the Steam Community Market page
		chromedp.Navigate("https://steamcommunity.com/market/listings/753/Frifle%20and%20Mauser"),

		// Wait for the page to load
		chromedp.WaitVisible(`#tabContentsMyActiveMarketListings`, chromedp.ByID),

		// Extract the price history variable
		chromedp.Evaluate(`JSON.stringify(line1);`, &priceHistory),
	)
	if err != nil {
		log.Fatalf("Failed to run chromedp tasks: %v", err)
	}

	// Parse the price history data
	var history [][]interface{}
	if err := json.Unmarshal([]byte(priceHistory), &history); err != nil {
		log.Fatalf("Failed to parse price history: %v", err)
	}

	// Print the price history
	fmt.Println("Price History:")
	for _, entry := range history {
		fmt.Printf("Date: %s, Price: %v, Volume: %v\n", entry[0], entry[1], entry[2])
	}

	fmt.Fprintf(w, "%s", history[0][1])
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

	// marketHashName := chi.URLParam(r, "market_hash_name")
	// if marketHashName == "" {
	// 	http.Error(w, "Market hash name is required", http.StatusBadRequest)
	// 	return
	// }
	// marketHashName = url.QueryEscape(marketHashName)
	// marketHashName = strings.ReplaceAll(marketHashName, "+", "%20")

	// Steam Market API URL for price history
	// url := fmt.Sprintf("https://steamcommunity.com/market/pricehistory/?appid=753&market_hash_name=%s", url.QueryEscape(marketHashName))

	//baseURL := fmt.Sprintf("https://steamcommunity.com/market/listings/753/%s/render/history", marketHashName)
	baseURL := "https://steamcommunity.com/market/itemordershistogram"
	u, err := url.Parse(baseURL)
	if err != nil {
		http.Error(w, "Failed to parse URL", http.StatusInternalServerError)
		return
	}

	q := u.Query()
	q.Add("country", "ES")
	q.Add("language", "english")
	q.Add("currency", "3") // 3 is for EUR
	q.Add("item_nameid", "150084592")
	u.RawQuery = q.Encode()

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
	} else {
		log.Printf("sessionid cookie not found")
	}

	// Add steamCountry cookie
	if steamCountry, ok := session.Values["steamCountry"].(string); ok {
		req.AddCookie(&http.Cookie{
			Name:  "steamCountry",
			Value: steamCountry,
		})
	} else {
		log.Printf("steamCountry cookie not found")
	}

	if steamLoginSecure, ok := session.Values["steamLoginSecure"].(string); ok {
		req.AddCookie(&http.Cookie{
			Name:  "steamLoginSecure",
			Value: steamLoginSecure,
		})
	} else {
		log.Printf("steamLoginSecure cookie not found")
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

	// log.Printf("Requesting URL: %s", u.String())
	// log.Printf("With cookies: %v", req.Cookies())

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to fetch market data", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// if resp.StatusCode != http.StatusOK {
	// 	log.Printf("Steam Error - Status: %d", resp.StatusCode)
	// 	log.Printf("Steam Error - Headers: %v", resp.Header)
	// 	log.Printf("Steam Error - URL: %s", baseURL)
	// }

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

	// log.Printf("Response Status: %d", resp.StatusCode)
	// log.Printf("Response Body: %s", string(body))

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
