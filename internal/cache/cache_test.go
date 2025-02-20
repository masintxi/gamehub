package cache

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/path",
			val: []byte("moretestdata"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := NewCache(CacheConfig{
				CleanupInterval: interval,
			})
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find key")
				return
			}
			if string(val) != string(c.val) {
				t.Errorf("expected to find value")
				return
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	const cleanupTime = 3 * time.Millisecond
	const expirationTime = 5 * time.Millisecond
	const waitTime = expirationTime + 5*time.Millisecond
	cache := NewCache(CacheConfig{
		CleanupInterval: cleanupTime,
		ExpireAfter:     expirationTime,
	})
	cache.Add("https://example.com", []byte("testdata"))

	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime)

	_, ok = cache.Get("https://example.com")
	if ok {
		t.Errorf("expected to not find key")
		return
	}
}

func TestMaxSize(t *testing.T) {
	cache := NewCache(CacheConfig{
		MaxSize: 10, // 10 bytes total
	})
	cache.ClearMemoryCache()

	key1Data := []byte("12345")
	key2Data := []byte("1234567")

	fmt.Printf("key1 size: %d bytes\n", len(key1Data))
	fmt.Printf("key2 size: %d bytes\n", len(key2Data))

	cache.Add("key0", []byte("123456789012345")) // 15 bytes
	_, ok := cache.Get("key0")
	if ok {
		t.Error("key0 should not have been added - too big")
	}

	cache.Add("key1", key1Data) // 5 bytes
	val1, ok := cache.Get("key1")
	if !ok || string(val1) != "12345" {
		t.Error("expected to find key1")
	}

	cache.Add("key2", key2Data) // 7 bytes

	_, ok = cache.Get("key1")
	if ok {
		t.Error("expected key1 to be evicted")
	}

	val2, ok := cache.Get("key2")
	if !ok || string(val2) != "1234567" {
		t.Error("expected to find key2")
	}

	cache.PrintCache()
}

func TestLoadCache(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cache-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	cache := NewCache(CacheConfig{
		MaxSize: 100,
	})
	cache.filePath = filepath.Join(tmpDir, "cache.json")

	cache.Add("key1", []byte("12345"))
	cache.Add("key2", []byte("67890"))

	expectedSize := int64(10) // 5 bytes + 5 bytes

	newCache := NewCache(CacheConfig{
		MaxSize: 100,
	})
	newCache.filePath = cache.filePath

	fmt.Print("Cache1 - ")
	cache.PrintCache()

	err = newCache.LoadCache()
	if err != nil {
		t.Fatalf("Failed to load cache: %v", err)
	}

	fmt.Print("Cache2 - ")
	newCache.PrintCache()

	if newCache.currentSize != expectedSize {
		t.Errorf("Expected size %d, got %d", expectedSize, newCache.currentSize)
	}
}
