package redis

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
)

/*
提供工厂函数，
提供获取锁（失败立即返回）接口，（锁续期功能）
提供获取锁（等待一段时间后超时立即返回），（锁续期功能）
*/

// NewRedisMutex 工厂函数，返回redis分布式锁句柄
func NewRedisMutex(addr, pwd, key string, timeout int) *RedisMutex {
	c := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd,
	})
	return &RedisMutex{
		client:  c,
		key:     key,
		timeout: timeout,
	}
}

// Key return redis key
func (m *RedisMutex) Key() string {
	return fmt.Sprintf("lock_%v", m.key)
}

// Lock 上锁
func (m *RedisMutex) Lock() error {
	res, err := safeLock.Run(context.Background(), m.client, []string{m.Key()}, m.timeout).Int64()
	if err != nil && err != redis.Nil {
		return err
	} else if res != 1 {
		return errors.New("redis transaction doesn't return 1")
	}
	return nil
}

// Lock 解锁
func (m *RedisMutex) Unlock() {

}
