package limiter

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

// args RedisKey 生成速度 容量 当前时间戳
var tokenBucketLua = `
local key = "rate:d:token:bucket:" .. KEYS[1]
local rate = tonumber(ARGV[1])
local cap = tonumber(ARGV[2])
local now = tonumber(ARGV[3])

local is_exists = redis.call("EXISTS", key)
if is_exists == 0 then
    redis.call("HMSET", key, "last_time", now, "tokens", cap, "cap", cap, "rate", rate)
end

local bucket = redis.pcall("HMGET", key, "last_time", "tokens", "cap", "rate")

local bucket_last_time = tonumber(bucket[1])

local bucket_tokens = tonumber(bucket[2])

local bucket_cap = tonumber(bucket[3])

local bucket_rate = tonumber(bucket[4])

if cap ~= bucket_cap or rate ~= bucket_rate then
    bucket_cap = cap
    bucket_rate = rate
    redis.pcall("HMSET", key, "cap", cap, "rate", rate)
end

local cur_tokens = bucket_tokens + (now-bucket_last_time)*bucket_rate

if cur_tokens > bucket_cap then
    cur_tokens = bucket_cap
end

local ret = 0
if cur_tokens > 0 then
    cur_tokens = cur_tokens - 1
    ret = 1
end

redis.pcall("HMSET", key, "last_time", now, "tokens", cur_tokens)
return ret
`

func NewDTokenBucket(redisKey string, rate, cap int, redisClient *redis.Client) *DTokenBucket {
	return &DTokenBucket{redisKey: redisKey, rate: rate, cap: cap, redisClient: redisClient}
}

type DTokenBucket struct {
	redisKey    string
	rate        int
	cap         int
	redisClient *redis.Client
}

func (dt *DTokenBucket) Take() bool {
	// 执行lua
	res, err := dt.redisClient.Eval(
		context.TODO(), tokenBucketLua, []string{dt.redisKey},
		dt.rate, dt.cap, time.Now().Unix(),
	).Result()
	if err != nil {
		return false
	}
	if code, ok := res.(int64); ok {
		return code == 1
	}
	return false
}
