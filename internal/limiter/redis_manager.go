package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var incrWithExpire = redis.NewScript(`
local count = redis.call("INCR", KEYS[1])
if count == 1 then
	redis.call("EXPIRE", KEYS[1], ARGV[1])
end
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
	
	// create unique key for user (ip) for current minute
	key := fmt.Sprintf("rate:%s:%s", ip, time.Now().Format("15:04"))

	// atomic increment + expire via lua script
	count, err := incrWithExpire.Run(ctx, rm.client, []string{key}, 60).Int64()
	if err != nil {
		fmt.Printf("Redis error: %v\n", err)
		return false // when redis fails, block traffic
	}

	return int(count) <= rm.limit
}