package redis

import (
	"github.com/go-redis/redis/v8"
)

type RedisMutex struct {
	client  *redis.Client
	key     string // not redis key, but lock tag
	timeout int
}

var safeLock = redis.NewScript(`
	local key = KEYS[1]
	local r = redis.call("SETNX", key, 1)
	if (r == 0) then
		return 0
	end

	redis.call("EXPIRE", key, ARGV[1])
	return 1
`)

var tmpTransaction = redis.NewScript(`
	redis.call("SET", "random", 1)
	redis.call("SET", "key", 2)
	return 1
`)
