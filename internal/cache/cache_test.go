package distributedcache

import (
	"sync"
	"testing"
	"time"
)

const CacheCleanupTick100ms = 100 * time.Millisecond

func Test_Get(t *testing.T) {
	cache := NewCache(CacheCleanupTick100ms)

	// Get Key: Key Not Found
	_, err := cache.Get([]byte("shouldErrkey"))

	if err == nil {
		t.Error("Should have errored out")
	}

	// Get Key: Key Found
	key := []byte("testKey")
	value := []byte("testVal")
	_ = cache.Set(key, value, 10)

	retrievedValue, err := cache.Get(key)

	if err != nil {
		t.Errorf("Error Get value from cache: %v", err)
	}

	if string(retrievedValue) != string(value) {
		t.Errorf("Expected: %s, Got: %s", value, retrievedValue)
	}
}

func Test_SetWithTTL_Cleanup(t *testing.T) {
	cache := NewCache(CacheCleanupTick100ms)
	ttl30ms := 30 * time.Millisecond

	cache.StartCleanup()
	defer cache.StopCleanup()

	key := []byte("keyWithTTL")
	value := []byte("value")

	err := cache.Set(key, value, ttl30ms)

	if err != nil {
		t.Fatalf("Set with TTL failed: %v", err)
	}

	if !cache.Has(key) {
		t.Fatalf("Expected key to exist right after SET, but no key: %s", key)
	}

	// Wait for the clean up to kick in
	time.Sleep(ttl30ms + CacheCleanupTick100ms + 20*time.Millisecond)

	if cache.Has(key) {
		t.Fatalf("Key %s should have been deleted but is not", key)
	}
}

func Test_Has(t *testing.T) {
	cache := NewCache(CacheCleanupTick100ms)

	key := []byte("testHasKey")
	value := []byte("testHasValue")

	_ = cache.Set(key, value, 10)

	if !cache.Has(key) {
		t.Error("Key should be present but is not")
	}

	if cache.Has([]byte("shouldNotExist")) {
		t.Error("This key should not exist but does")
	}
}

func Test_Delete(t *testing.T) {
	cache := NewCache(CacheCleanupTick100ms)

	key := []byte("testDelete")

	_ = cache.Set(key, []byte("testDeleteValue"), 10)

	_ = cache.Delete(key)

	if cache.Has(key) {
		t.Error("Deletion failed")
	}

	// Attempt to delete a non-existent key
	nonExistentKey := []byte("nonExistentKey")
	err := cache.Delete(nonExistentKey)

	if err == nil {
		t.Error("shouldn't be able to delete nonExistentKey")
	}
}

func Test_SetWithZeroTTL(t *testing.T) {
	cache := NewCache(CacheCleanupTick100ms)

	key := []byte("testZeroTTLKey")
	value := []byte("testZeroTTLValue")
	err := cache.Set(key, value, 10)

	if err != nil {
		t.Errorf("Expected to Set key: %s, value: %s, Got: %v", key, value, err)
	}

	if !cache.Has(key) {
		t.Error("Expected key to persist with zero TTL but it does not")
	}
}

func Test_SetWithNegativeTTL(t *testing.T) {
	cache := NewCache(CacheCleanupTick100ms)

	key := []byte("testNegativeTTLKey")
	value := []byte("testNegativeTTLValue")
	err := cache.Set(key, value, -time.Second)

	if err == nil {
		t.Errorf("Expected an error when setting key: %s with negative TTL, but got none", key)
	}
}

func Test_ConcurrentAccess(t *testing.T) {
	cache := NewCache(CacheCleanupTick100ms)
	key := []byte("testConcurrentKey")
	value := []byte("testConcurrentValue")

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// nolint:errcheck
			cache.Set(key, value, 10)
			if !cache.Has(key) {
				t.Error("Expected key to be present but it is not")
			}
			retrievedValue, err := cache.Get(key)
			if err != nil {
				t.Errorf("Error getting value from cache: %v", err)
			}
			if string(retrievedValue) != string(value) {
				t.Errorf("Expected: %s, Got: %s", value, retrievedValue)
			}
		}()
	}
	wg.Wait()
}
