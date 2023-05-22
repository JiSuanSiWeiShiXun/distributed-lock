package redis

import (
	"github.com/go-redis/redis/v8"
)

type RedisMutex struct {
	client  *redis.Client
	key     string // not redis key, but lock tag
	timeout int
}

// redis获取锁操作的lua实现
var safeLock = redis.NewScript(`
	local key = KEYS[1]
	local r = redis.call("SETNX", key, 1)
	if (r == 0) then
		return 0
	end

	redis.call("EXPIRE", key, ARGV[1])
	return 1
`)

// redis释放锁
var unlock = redis.NewScript(`
	redis.call("del", KEYS[1])
	return "OK"
`)

// 测试发现有个坑，这个示例中的两个set命令即便执行成功，这个判断语句也不会走进return "OK"，而是会走进return "SHIT"
var tmpTransaction = redis.NewScript(`
	local r1 = redis.call('SET', KEYS[1], ARGV[1])
	local r2 = redis.call('SET', KEYS[2], ARGV[2])

	--if r1 == "OK" and r2 == "OK" then
	-- 	return "OK"
	--else
	--	return "SHIT"
	--end

	return "OK"
`)
