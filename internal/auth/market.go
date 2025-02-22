package auth

import (
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
)

func (sa *SteamAuth) InitializeMarketSession() error {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return fmt.Errorf("failed to create cookie jar: %w", err)
	}

	client := &http.Client{
		Jar: jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil
		},
	}

	req, err := http.NewRequest("GET", "https://steamcommunity.com/market/", nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to initialize market session: %w", err)
	}
	defer resp.Body.Close()

	// Extract cookies from response
	for _, cookie := range jar.Cookies(resp.Request.URL) {
		switch cookie.Name {
		case "sessionid":
			sa.SessionID = cookie.Value
			log.Printf("Got sessionid: %s", cookie.Value)
		case "steamLoginSecure":
			sa.SteamLoginSecure = cookie.Value
			log.Printf("Got steamLoginSecure: %s", cookie.Value)
		}
	}

	if sa.SessionID == "" || sa.SteamLoginSecure == "" {
		return fmt.Errorf("failed to capture required Steam cookies")
	}

	return nil
}
