package store

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type RateLimiterPersistentStore struct {
	visitors map[string]*Visitor
	mutex    sync.Mutex
	rate     rate.Limit // for more info check out Limiter docs - https://pkg.go.dev/golang.org/x/time/rate#Limit.

	burst       int
	expiresIn   time.Duration
	lastCleanup time.Time

	timeNow func() time.Time
}

// Visitor signifies a unique user's limiter details
type Visitor struct {
	*rate.Limiter
	lastSeen time.Time
}

/*
cleanupStaleVisitors helps manage the size of the visitors map by removing stale records
of users who haven't visited again after the configured expiry time has elapsed
*/
func (store *RateLimiterPersistentStore) cleanupStaleVisitors() {
	for id, visitor := range store.visitors {
		if store.timeNow().Sub(visitor.lastSeen) > store.expiresIn {
			delete(store.visitors, id)
		}
	}
	store.lastCleanup = store.timeNow()
}

func NewRateLimiterPersistentStoreWithConfig(config RateLimiterPersistentStoreConfig) (store *RateLimiterPersistentStore) {
	store = &RateLimiterPersistentStore{}

	store.rate = config.Rate
	store.burst = config.Burst
	store.expiresIn = config.ExpiresIn
	if config.ExpiresIn == 0 {
		store.expiresIn = DefaultRateLimiterMemoryStoreConfig.ExpiresIn
	}
	if config.Burst == 0 {
		store.burst = int(config.Rate)
	}
	store.visitors = make(map[string]*Visitor)
	store.timeNow = time.Now
	store.lastCleanup = store.timeNow()
	return
}

// RateLimiterPersistentStoreConfig represents configuration for RateLimiterMemoryStore
type RateLimiterPersistentStoreConfig struct {
	Rate      rate.Limit    // Rate of requests allowed to pass as req/s. For more info check out Limiter docs - https://pkg.go.dev/golang.org/x/time/rate#Limit.
	Burst     int           // Burst is maximum number of requests to pass at the same moment. It additionally allows a number of requests to pass when rate limit is reached.
	ExpiresIn time.Duration // ExpiresIn is the duration after that a rate limiter is cleaned up
}

// DefaultRateLimiterMemoryStoreConfig provides default configuration values for RateLimiterMemoryStore
var DefaultRateLimiterMemoryStoreConfig = RateLimiterPersistentStoreConfig{
	ExpiresIn: 3 * time.Minute,
}

func (store *RateLimiterPersistentStore) Allow(identifier string) (bool, error) {
	store.mutex.Lock()
	limiter, exists := store.visitors[identifier]
	if !exists {
		limiter = new(Visitor)
		limiter.Limiter = rate.NewLimiter(store.rate, store.burst)
		store.visitors[identifier] = limiter
	}
	now := store.timeNow()
	limiter.lastSeen = now
	if now.Sub(store.lastCleanup) > store.expiresIn {
		store.cleanupStaleVisitors()
	}
	store.mutex.Unlock()
	return limiter.AllowN(store.timeNow(), 1), nil
}
