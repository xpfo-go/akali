package limiter

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

var leakyBucketLua = `
local key = "rate:d:leaky:bucket:" .. KEYS[1]
local peak = tonumber(ARGV[1])
local rate = tonumber(ARGV[2])
local now = tonumber(ARGV[3])

local is_exists = redis.call("EXISTS", key)
if is_exists == 0 then
    redis.call("HMSET", key, "last_time", now, "peak", peak, "rate", rate, "cur", 0)
end

local bucket = redis.pcall("HMGET", key, "last_time", "peak", "rate", "cur")
local bucket_last_time = tonumber(bucket[1])
local bucket_peak = tonumber(bucket[2])
local bucket_rate = tonumber(bucket[3])
local bucket_cur = tonumber(bucket[4])

if rate ~= bucket_rate or peak ~= bucket_peak then
    bucket_rate = rate
    bucket_peak = peak
    redis.pcall("HMSET", key, "peak", peak, "rate", rate)
end

local cur = bucket_cur - (now-bucket_last_time)*bucket_rate
if cur {{ .Lt }} 0 then
    cur = 0
end

if cur >= bucket_peak then
    redis.pcall("HMSET", key, "last_time", now)
    return 0
end

redis.pcall("HMSET", key, "last_time", now, "cur", cur+1)

return 1
`

func NewDLeakyBucket(redisKey string, rate, peak int, redisClient *redis.Client) *DLeakyBucket {
	return &DLeakyBucket{redisKey: redisKey, rate: rate, peak: peak, redisClient: redisClient}
}

type DLeakyBucket struct {
	redisKey    string
	rate        int
	peak        int
	redisClient *redis.Client
}

func (dl *DLeakyBucket) Take() bool {
	// 执行lua
	res, err := dl.redisClient.Eval(context.TODO(), leakyBucketLua, []string{dl.redisKey}, dl.peak, dl.rate, time.Now().Unix()).Result()
	if err != nil {
		return false
	}

	if code, ok := res.(int64); ok {
		return code == 1
	}
	return false
}
