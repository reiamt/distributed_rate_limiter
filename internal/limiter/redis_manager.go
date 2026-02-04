package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisManager struct {
	client	*redis.Client
	limit	int
}

func NewRedisManager(addr string, limit int) *RedisManager {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	return &RedisManager{
		client:	rdb,
		limit:	limit,
	}
}

func (rm *RedisManager) Allow(ip string) bool {
	ctx := context.Background()
	
	// create unique key for user (ip) for current minute
	key := fmt.Sprintf("rate:%s:%s", ip, time.Now().Format("15:04"))

	// atomic increment
	count, err := rm.client.Incr(ctx, key).Result()
	if err != nil {
		fmt.Printf("Redis error: %v\n", err)
		return false // when redis fail, allow traffic
	}

	// if its new key, set it to expire
	if count == 1 {
		rm.client.Expire(ctx, key, time.Minute)
	}

	return int(count) <= rm.limit
}