package client

import (
	"net/http"
	"time"

	"github.com/masintxi/gamehub/internal/cache"
)

type Client struct {
	HttpClient *http.Client
	Cache      *cache.Cache
}

func NewClient(cacheConfig cache.CacheConfig) *Client {
	return &Client{
		HttpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		Cache: cache.NewCache(cacheConfig),
	}
}
