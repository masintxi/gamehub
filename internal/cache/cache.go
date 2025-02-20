package cache

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type cacheEntry struct {
	CreatedAt time.Time `json:"created_at"`
	Val       []byte    `json:"val"`
	Size      int64     `json:"size"`
}

type CacheConfig struct {
	ProjectName     string
	CleanupInterval time.Duration
	MaxSize         int64         // Maximum size in MB
	FileExtension   string        // "json", "gob", etc
	Compression     bool          // Whether to use gzip
	ExpireAfter     time.Duration // How long items stay valid
	CachePath       string        // Optional custom path override
}

type Cache struct {
	cache       map[string]cacheEntry
	mu          *sync.Mutex
	filePath    string
	config      CacheConfig
	currentSize int64
}

func NewCache(config CacheConfig) *Cache {
	if config.ProjectName == "" {
		config.ProjectName = "unnamed-project"
	}
	if config.CleanupInterval == 0 {
		config.CleanupInterval = 30 * time.Minute
	}
	if config.MaxSize == 0 {
		config.MaxSize = 1024 * 1024 * 10 // 10 MB
	}
	if config.FileExtension == "" {
		config.FileExtension = "json"
	}
	if config.ExpireAfter == 0 {
		config.ExpireAfter = 30 * time.Minute
	}

	c := &Cache{
		cache:  make(map[string]cacheEntry),
		mu:     &sync.Mutex{},
		config: config,
	}
	if err := c.CreateCacheDir(); err != nil {
		log.Printf("Warning: Cache directory creation failed: %v. Continuing with in-memory cache only\n", err)
	}

	if err := c.LoadCache(); err != nil {
		log.Printf("Error loading cache: %v\n", err)
	}

	go c.reapLoop(config.CleanupInterval)

	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	newSize := int64(len(val))

	if newSize > c.config.MaxSize {
		log.Printf("Warning: Item size %d bytes exceeds cache max size %d bytes",
			newSize, c.config.MaxSize)
		return
	}

	for c.currentSize+newSize > c.config.MaxSize {
		var oldestKey string
		var oldestTime time.Time

		for k, v := range c.cache {
			//log.Printf("  Key: %s, Created: %v", k, v.CreatedAt)
			if oldestKey == "" || v.CreatedAt.Before(oldestTime) {
				oldestKey = k
				oldestTime = v.CreatedAt
			}
		}

		if oldestKey != "" {
			removedSize := c.cache[oldestKey].Size
			delete(c.cache, oldestKey)
			c.currentSize -= removedSize
			//log.Printf("Removed cache entry %s to free %d bytes", oldestKey, removedSize)
		}
	}

	c.cache[key] = cacheEntry{
		CreatedAt: time.Now(),
		Val:       val,
		Size:      newSize,
	}
	c.currentSize += newSize

	if err := c.SaveCache(); err != nil {
		log.Printf("Error saving cache: %v", err)
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if entry, ok := c.cache[key]; ok {
		return entry.Val, true
	}
	return nil, false
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, entry := range c.cache {
			if now.Sub(entry.CreatedAt) > c.config.ExpireAfter {
				c.currentSize -= entry.Size
				delete(c.cache, key)
			}
		}
		c.mu.Unlock()
	}
}

func (c *Cache) PrintCache() {
	c.mu.Lock()
	defer c.mu.Unlock()
	fmt.Println("Cache Contents:")
	for key, entry := range c.cache {
		fmt.Printf("Key: %s, Value: %v\n", key, entry.Val)
	}
}

func (c *Cache) SaveCache() error {
	data, err := json.Marshal(c.cache)
	if err != nil {
		return fmt.Errorf("marshalling cache: %w", err)
	}

	file, err := os.Create(c.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if c.config.Compression {
		gw := gzip.NewWriter(file)
		defer gw.Close()

		_, err = gw.Write(data)
		return err

	} else {
		_, err = file.Write(data)
		return err
	}
}

func (c *Cache) LoadCache() error {
	file, err := os.Open(c.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	var reader io.Reader = file
	if c.config.Compression {
		gr, err := gzip.NewReader(file)
		if err != nil {
			return err
		}
		defer gr.Close()
		reader = gr
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	c.currentSize = 0

	if err := json.Unmarshal(data, &c.cache); err != nil {
		return err
	}

	for _, entry := range c.cache {
		c.currentSize += entry.Size
	}

	return nil
}

func (c *Cache) CreateCacheDir() error {
	c.filePath = getCacheFilePath(c.config)
	return os.MkdirAll(filepath.Dir(c.filePath), 0755)
}

func (c *Cache) DeleteCacheDir() error {
	c.ClearMemoryCache()
	return os.RemoveAll(filepath.Dir(c.filePath))
}

func (c *Cache) ClearMemoryCache() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache = make(map[string]cacheEntry)
	c.currentSize = 0
}

func (c *Cache) GetCacheSize() (float64, error) {
	fileInfo, err := os.Stat(c.filePath)
	if err != nil {
		return 0, fmt.Errorf("error getting cache file info: %w", err)
	}

	sizeMB := float64(fileInfo.Size()) / 1024 / 1024

	return sizeMB, nil
}

func getCacheFilePath(config CacheConfig) string {
	basePath := config.CachePath
	if basePath == "" {
		if path, err := os.UserCacheDir(); err == nil {
			basePath = path
		} else if path, err := os.UserHomeDir(); err == nil {
			basePath = path
		} else {
			tmpDir := os.TempDir()
			if _, err := os.Stat(tmpDir); err != nil {
				log.Printf("Warning: default temp directory %s is not accessible: %v\n", tmpDir, err)
				log.Println("Falling back to current directory")
				basePath = "."
			} else {
				basePath = tmpDir
			}
		}
	}

	cacheDir := filepath.Join(basePath, config.ProjectName+"-cache")

	filename := "cache"
	if config.Compression {
		filename = filename + ".gz"
	} else {
		filename = filename + "." + config.FileExtension
	}

	return filepath.Join(cacheDir, filename)
}
