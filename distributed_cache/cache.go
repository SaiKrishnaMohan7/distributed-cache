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

type InMemoryCache struct {
	data map[string][]byte
	expiry map[string][]byte
	lock sync.RWMutex // Concurrent read/write safety
	cleanupTick time.Duration
	stopCleanup chan struct{}
}

func NewCache() *InMemoryCache {
	// This is called a constructor fn, idiomatic way
	// to initialize a struct and encapsulate info
	return &InMemoryCache{
		data: make(map[string][]byte),
	}
}

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

func (cache *InMemoryCache) Set(key, value []byte, ttl time.Duration) error {
	if ttl < 0 {
		return fmt.Errorf("ttl cannot be negative")
	}

	cache.lock.Lock()
	defer cache.lock.Unlock()

	keyStr := string(key)

	cache.data[keyStr] = value

	if ttl > 0 {
		// TODO: Do not spawn a go routine for every Set operation, wasteful
		go func() {
			// time.After() returns a buffered RECEIVE ONLY timerChannel
			// A value is received by this channel when ttl expires
			// the value received is timer.Time type
			timerChannel := time.After(ttl)
			// the received value is discarded as there is no use for it
			// <- operator here, blocks the channel till a value is sent from it!
			<-timerChannel
			cache.lock.Lock()
			defer cache.lock.Unlock()
			delete(cache.data, keyStr)
		}()
	}

	return nil
}

func (cache *InMemoryCache) Has(key []byte) bool {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	_, ok := cache.data[string(key)]

	return ok
}

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
