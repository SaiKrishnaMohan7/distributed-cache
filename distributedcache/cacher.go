package distributedcache

import (
	"fmt"
	"sync"
	"time"
)

type Cache interface {
	Set([]byte, []byte, time.Duration) error
	Get([]byte) ([]byte, error)
	Delete([]byte) error
	Has([]byte) bool
}

type InMemoryCache struct {
	lock sync.RWMutex // Concurrent read/write safety

	data map[string][]byte
}

func New() *InMemoryCache {
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

	delete(cache.data, string(key))

	return nil
}
