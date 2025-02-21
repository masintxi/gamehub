package handlers

import (
	"fmt"
	"io"
	"net/http"
)

func (h *SteamHandlers) getResponseBody(url string, headers map[string]string) ([]byte, error) {
	// Try to get the data from the cache first
	val, ok := h.client.Cache.Get(url)
	if ok {
		return val, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	resp, err := h.client.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	h.client.Cache.Add(url, []byte(string(body)))

	return body, nil
}
