package limiter

import (
	"testing"
	"time"
)

func TestNewBucketAllowed(t *testing.T) {
	capacity := 5
	tokenBucket := NewTokenBucket(float64(capacity), 1.0)

	// all succeed
	for cntr := 0; cntr < capacity; cntr++ {
		allowed := tokenBucket.Allow().Allowed
		if !allowed {
			t.Errorf("Request should be allowed.")
		}
	}
	
	// capacity+1-th request fails
	if tokenBucket.Allow().Allowed {
		t.Errorf("Request should not be allowed.")
	}
}

func TestNoTokensFail(t *testing.T) {
	capacity := 0
	tokenBucket := NewTokenBucket(float64(capacity), 1.0)

	// request fails, no tokens
	if tokenBucket.Allow().Allowed {
		t.Errorf("Request should not be allowed.")
	}	
}

func TestTokenRefill(t *testing.T) {
	capacity := 3
	tokenBucket := NewTokenBucket(float64(capacity), 10.0)

	// use all tokens
	for cntr := 0; cntr < capacity; cntr++ {
		allowed := tokenBucket.Allow().Allowed
		if !allowed {
			t.Errorf("Request should be allowed.")
		}
	}

	// refill 1 token
	time.Sleep(100 * time.Millisecond)
	allowed := tokenBucket.Allow().Allowed
	if !allowed {
		t.Errorf("Request should be allowed.")
	}
}

func TestBucketCapacityCeil(t *testing.T) {
	capacity := 2
	tokenBucket := NewTokenBucket(float64(capacity), 10.0)

	// if bucket not full, would add 2 tokens
	time.Sleep(2 * 100 * time.Millisecond)

	for cntr := 0; cntr < capacity; cntr++ {
		allowed := tokenBucket.Allow().Allowed
		if !allowed {
			t.Errorf("Request should be allowed.")
		}
	}

	if tokenBucket.Allow().Allowed {
		t.Errorf("Request should not be allowed.")
	}
}

func TestGetBucketSameKey(t *testing.T) {
	mgr := NewManager(3.0, 1.0)
	
	bucket1 := mgr.GetBucket("ip1")
	bucket2 := mgr.GetBucket("ip1")

	if bucket1 != bucket2 {
		t.Errorf("Pointers should be the same.")
	}
}

func TestGetBucketDifferentKey(t *testing.T) {
	mgr := NewManager(3.0, 1.0)

	bucket1 := mgr.GetBucket("ip1")
	bucket2 := mgr.GetBucket("ip2")

	if bucket1 == bucket2 {
		t.Errorf("Pointers should not be the same.")
	}
}

func TestPerKeyLimit(t *testing.T) {
	capacity := 3
	mgr := NewManager(float64(capacity), 1.0)

	bucket1 := mgr.GetBucket("ip1")
	bucket2 := mgr.GetBucket("ip2")

	for cntr := 0; cntr < capacity; cntr++ {
		allowed := bucket1.Allow().Allowed
		if !allowed {
			t.Errorf("Bucket1 request should be allowed.")
		}
	}

	if bucket1.Allow().Allowed {
		t.Errorf("Bucket1 request should not be allowed.")
	}

	if !bucket2.Allow().Allowed {
		t.Errorf("Bucket2 request should be allowed.")
	}
}