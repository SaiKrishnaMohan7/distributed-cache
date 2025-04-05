package distributedcache

import (
	"fmt"
	"sync"
	"time"
)

// Cache is an interface that defines the basic operations for a distributed cache.
// It includes methods for setting, getting, deleting, and checking the existence of keys.
type Cache interface {
	Set([]byte, []byte, time.Duration) error
	Get([]byte) ([]byte, error)
	Delete([]byte) error
	Has([]byte) bool
}

// InMemoryCache is an in-memory implementation of the Cache interface.
// It provides thread-safe operations and supports key expiration with cleanup.
type InMemoryCache struct {
	data map[string][]byte
	expiry map[string]time.Time
	lock sync.RWMutex // Concurrent read/write safety
	cleanupTick time.Duration
	stopCleanup chan struct{}
}

// NewCache creates a new instance of InMemoryCache with a specified cleanup interval.
// The cleanupTick parameter defines the duration between periodic cleanup operations
// to remove expired items from the cache.
func NewCache(cleanupTick time.Duration) *InMemoryCache {
	// This is called a constructor fn, idiomatic way
	// to initialize a struct and encapsulate info
	return &InMemoryCache{
		data: make(map[string][]byte),
		expiry: make(map[string]time.Time),
		cleanupTick: cleanupTick,
		stopCleanup: make(chan struct{}),
	}
}

// Get retrieves the value associated with the given key from the in-memory cache.
// If the key exists, it returns the corresponding value and a nil error.
// If the key does not exist, it returns a nil value and an error.
func (cache *InMemoryCache) Get(key []byte) ([]byte, error) {
	// Acquire a read lock
	cache.lock.RLock()
	defer cache.lock.RUnlock()

	keyStr := string(key)

	value, ok := cache.data[keyStr]

	if !ok {
		return nil, fmt.Errorf("key (%s) not found", keyStr)
	}

	return value, nil

}

// Set stores a key-value pair in the cache with an optional time-to-live (TTL).
// If the TTL is greater than zero, the key will expire after the specified duration.
// If the TTL is zero, the key will persist indefinitely.
func (cache *InMemoryCache) Set(key, value []byte, ttl time.Duration) error {
	if ttl < 0 {
		return fmt.Errorf("ttl cannot be negative")
	}

	cache.lock.Lock()
	defer cache.lock.Unlock()

	keyStr := string(key)

	cache.data[keyStr] = value

	if ttl > 0 {
		cache.expiry[keyStr] = time.Now().Add(ttl)
	} else {
		delete(cache.expiry, keyStr)
	}

	return nil
}

// Has checks if a given key exists in the cache.
// It returns true if the key is present, otherwise false.
func (cache *InMemoryCache) Has(key []byte) bool {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	_, ok := cache.data[string(key)]

	return ok
}

// Delete removes the specified key from the cache.
// If the key exists, it is deleted, and nil is returned.
// If the key does not exist, an error is returned.
func (cache *InMemoryCache) Delete(key []byte) error {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	_, ok := cache.data[string(key)]

	if ok {
		delete(cache.data, string(key))
		return nil
	}

	return fmt.Errorf("invalid key: %s", key)
}

// StartCleanup initiates a periodic cleanup process to remove expired keys from the cache.
// It runs in a separate goroutine and stops when StopCleanup is called.
func (cache *InMemoryCache) StartCleanup() {
	ticker := time.NewTicker(cache.cleanupTick)

	go func () {
		for {
			select {
			case <-ticker.C:
				now := time.Now()

				if (len(cache.expiry) != 0) {
					// nothing to clean, no work needed
					cache.lock.Lock()
				}

				for key, expiry := range cache.expiry {
					if (now.After(expiry)) {
						delete(cache.data, key)
						delete(cache.expiry, key)
					}
				}
				// we want to release lock BEFORE the goroutine is done.
				// ensures locks are not held any longer than they have to
				cache.lock.Unlock()
			case <-cache.stopCleanup:
				ticker.Stop()
				return
			}
		}
	}()
}

// StopCleanup stops the periodic cleanup process by closing the stopCleanup channel.
// This ensures that the cleanup goroutine terminates gracefully.
func (cache *InMemoryCache) StopCleanup() {
	close(cache.stopCleanup)
}