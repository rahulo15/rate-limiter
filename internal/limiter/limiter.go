package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

// Lua Script (The Logic)
var tokenBucketScript = redis.NewScript(`
	local key_tokens = KEYS[1]
	local key_time = KEYS[2]
	local capacity = tonumber(ARGV[1])
	local rate = tonumber(ARGV[2])
	local now = tonumber(ARGV[3])
	local requested = 1

	local last_tokens = tonumber(redis.call("GET", key_tokens))
	local last_time = tonumber(redis.call("GET", key_time))

	if last_tokens == nil then
		last_tokens = capacity
		last_time = now
	end

	local delta = math.max(0, now - last_time)
	local filled_tokens = math.min(capacity, last_tokens + (delta * rate))

	local allowed = 0
	local new_tokens = filled_tokens

	if filled_tokens >= requested then
		allowed = 1
		new_tokens = filled_tokens - requested
	end

	redis.call("SET", key_tokens, new_tokens)
	redis.call("SET", key_time, now)

	return allowed
`)

func AllowRequest(rdb *redis.Client, userID string, capacity int, rate float64) bool {
	tokensKey := fmt.Sprintf("rate_limit:%s:tokens", userID)
	tsKey := fmt.Sprintf("rate_limit:%s:ts", userID)

	now := time.Now().Unix()

	result, err := tokenBucketScript.Run(ctx, rdb, []string{tokensKey, tsKey}, capacity, rate, now).Result()
	if err != nil {
		fmt.Printf("Redis Error: %v\n", err)
		return true
	}

	return result.(int64) == 1
}
