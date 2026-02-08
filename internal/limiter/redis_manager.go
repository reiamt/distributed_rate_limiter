package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

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
	key := fmt.Sprintf("rate:%s", ip)
	now := time.Now().UnixMicro()
	window := int64(60 * 1e6) // 60 secs in microsecs

	// atomic increment + expire via lua script
	count, err := slidingWindowScript.Run(ctx, rm.client, []string{key}, window, now).Int64() 
	if err != nil {
		fmt.Printf("Redis error: %v\n", err)
		return false // when redis fails, block traffic
	}

	return int(count) <= rm.limit
}