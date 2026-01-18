package limiter

import (
	"sync"
	"time"
)

type TokenBucket struct {
	capacity	float64		// max tokens the bucket can hold
	rate		float64		// tokens added per sec
	tokens		float64		// current available tokens
	lastTick	time.Time	// last update
	mu			sync.Mutex
}

func NewTokenBucket(capacity, rate float64) *TokenBucket {
	return &TokenBucket{
		capacity:	capacity,
		rate: 		rate,
		tokens: 	capacity,
		lastTick:	time.Now(),
	}
}

func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastTick).Seconds()
	tokensToAdd := elapsed * tb.rate

	tb.tokens += tokensToAdd
	if tb.tokens > tb.capacity {
		tb.tokens = tb.capacity
	}
	tb.lastTick = now

	if tb.tokens > 1 {
		tb.tokens--
		return true
	}

	return false
}

type Manager struct {
	mu		sync.RWMutex
	buckets	map[string]*TokenBucket
	rate	float64
	capacity float64
}

func NewManager(capacity, rate float64) *Manager {
	return &Manager{
		buckets:	make(map[string]*TokenBucket),
		rate:		rate,
		capacity:	capacity,
	}
}

func (m *Manager) GetBucket(key string) *TokenBucket {
	// readlock; check if bucket already exists
	m.mu.RLock()
	bucket, exists := m.buckets[key]
	m.mu.RUnlock()

	if exists {
		return bucket
	}

	// writelock; create new bucket
	m.mu.Lock()
	defer m.mu.Unlock()

	if bucket, exists = m.buckets[key]; exists {
		return bucket
	}

	newBucket :=  NewTokenBucket(m.capacity, m.rate)
	m.bucket[key] = newBucket
	return newBucket
}