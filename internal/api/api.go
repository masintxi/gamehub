package api

// import (
// 	"log"
// 	"net/http"
// 	"os"

// 	"github.com/joho/godotenv"
// 	"github.com/masintxi/gamehub/internal/cache"
// )

// type Client struct {
// 	httpClient *http.Client
// 	cache      *cache.Cache
// 	apiKey     string
// 	steamID    string
// }

// func NewClient(cacheConfig cache.CacheConfig) *Client {
// 	godotenv.Load()
// 	apiKey := os.Getenv("STEAM_API_KEY")
// 	if apiKey == "" {
// 		log.Fatal("STEAM_API_KEY not setted. Check .env file")
// 	}
// 	steamID := os.Getenv("STEAM_USER_ID")
// 	if steamID == "" {
// 		log.Fatal("STEAM_USER_ID not setted. Check .env file")
// 	}
// 	// jar, _ := cookiejar.New(nil)

// 	return &Client{
// 		httpClient: &http.Client{
// 			// Jar: jar,
// 		},
// 		apiKey:  apiKey,
// 		steamID: steamID,
// 		cache:   cache.NewCache(cacheConfig),
// 	}
// }
