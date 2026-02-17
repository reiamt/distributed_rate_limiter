package limiter

import (
	"context"
	"distributed_rate_limiter/internal/metrics"

	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

// lua script for sliding window
var slidingWindowScript = redis.NewScript(`
local key = KEYS[1]
local window = tonumber(ARGV[1])
local now = tonumber(ARGV[2])

redis.call("ZREMRANGEBYSCORE", key, 0, now - window)
redis.call("ZADD", key, now, tostring(now))
local count = redis.call("ZCARD", key)
redis.call("EXPIRE", key, math.ceil(window / 1000000))

return count
`)

type RedisManager struct {
	client   *redis.Client
	limit    int
	cb       *CircuitBreaker
	failOpen bool
}

func NewRedisManager(addr string, limit int, failOpen bool) *RedisManager {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	cb := NewCircuitBreaker(3, 30*time.Second)

	return &RedisManager{
		client:   rdb,
		limit:    limit,
		cb:       cb,
		failOpen: failOpen,
	}
}

func (rm *RedisManager) Allow(ip string) Result {
	ctx := context.Background()
	key := fmt.Sprintf("rate:%s", ip)
	now := time.Now().UnixMicro()
	window := int64(60 * 1e6) // 60 secs in microsecs

	// check circuitbreaker before calling redis
	if !rm.cb.Allow() {
		return Result{
			Allowed:   rm.failOpen,
			Limit:     rm.limit,
			Remaining: max(rm.limit-1, 0),
			ResetAt:   time.Now().Unix() + 60,
		}
	}

	// atomic increment + expire via lua script
	count, err := slidingWindowScript.Run(ctx, rm.client, []string{key}, window, now).Int64()
	if err != nil {
		metrics.RedisErrorsTotal.Inc()
		rm.cb.RecordResult(false)
		slog.Error("redis error", "err", err)
		return Result{
			Allowed:   false,
			Limit:     rm.limit,
			Remaining: max(rm.limit-int(count), 0),
			ResetAt:   time.Now().Unix() + 60,
		} // when redis fails, block traffic
	}
	rm.cb.RecordResult(true)
	return Result{
		Allowed:   int(count) <= rm.limit,
		Limit:     rm.limit,
		Remaining: max(rm.limit-int(count), 0),
		ResetAt:   time.Now().Unix() + 60,
	}
}

func (rm *RedisManager) Close() error {
	return rm.client.Close()
}

func (rm *RedisManager) Ping() bool {
	ctx := context.Background()
	pong, _ := rm.client.Ping(ctx).Result()
	if pong != "PONG" {
		return false
	}
	return true
}
