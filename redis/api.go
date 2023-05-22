package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

/*
提供工厂函数，
提供获取锁（失败立即返回）接口，（锁续期功能）
提供获取锁（等待一段时间后超时立即返回），（锁续期功能）
*/

// NewRedisMutex 工厂函数，返回redis分布式锁句柄
// timeout 单位为秒
func NewRedisMutex(c *redis.Client, key string, timeout int) (*RedisMutex, error) {
	if res, err := c.Ping(context.Background()).Result(); err != nil || res != "PONG" {
		return nil, fmt.Errorf("invalid redis client, ping returns %v", res)
	}
	return &RedisMutex{
		client:  c,
		key:     key,
		timeout: timeout,
	}, nil
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
		return fmt.Errorf("redis transaction failed: %v", res)
	}
	return nil
}

// Lock 解锁
func (m *RedisMutex) Unlock() error {
	res, err := unlock.Run(context.Background(), m.client, []string{m.Key()}).Int64()
	if err != nil && err != redis.Nil {
		return err
	} else if res != 1 {
		return fmt.Errorf("redis transaction failed: %v", res)
	}
	return nil
}

// LockWithTimeout 尝试加锁，失败重试，直到超时
func (m *RedisMutex) LockWithTimeout(timeout time.Duration) error {
	timer := time.NewTimer(timeout)
	tickr := time.NewTicker(100 * time.Millisecond)
	for {
		select {
		case <-tickr.C:
			if err := m.Lock(); err != nil {
				break
			} else {
				return nil
			}
		case <-timer.C:
			return fmt.Errorf("[key] %v unavailable", m.Key())
		}
	}
}
