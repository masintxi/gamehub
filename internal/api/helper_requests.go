package api

import (
	"fmt"
	"io"
	"net/http"
)

func (client *Client) getResposeBody(url string) ([]byte, error) {
	// Try to get the data from the cache first
	val, ok := client.cache.Get(url)
	if ok {
		return val, nil
	}

	resp, err := http.Get(url)
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

	client.cache.Add(url, []byte(string(body)))

	return body, nil
}
