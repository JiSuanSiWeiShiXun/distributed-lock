package redis

import (
	"sync"
	"testing"
	"time"
)

func TestNewRedisMetex(t *testing.T) {
	setUp()
	_, err := NewRedisMutex(red, "random", 1) // 过期
	if err != nil {
		panic(err)
	}
}

func TestLock(t *testing.T) {
	setUp()
	wg := sync.WaitGroup{}
	n := 10
	wg.Add(n)
	for i := 0; i < n; i++ {
		time.Sleep(150 * time.Millisecond)
		go func() {
			defer wg.Done()
			redLock, err := NewRedisMutex(red, "random", 1) //
			if err != nil {
				panic(err)
			}
			if err := redLock.Lock(); err != nil {
				t.Errorf("get lock [%v] failed %v\n", redLock.Key(), err)
				return
			}
			t.Logf("get lock")
		}()
	}
	wg.Wait()
}

func TestLockWithTimeout(t *testing.T) {
	setUp()
	wg := sync.WaitGroup{}
	n := 10
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			redLock, err := NewRedisMutex(red, "random", 1) //
			if err != nil {
				panic(err)
			}
			if err := redLock.LockWithTimeout(10 * time.Second); err != nil {
				t.Errorf("get lock [%v] failed %v\n", redLock.Key(), err)
				return
			}
			t.Logf("get lock")
		}()
	}
	wg.Wait()
}

func TestUnlock(t *testing.T) {
	setUp()

	wg := sync.WaitGroup{}
	n := 10
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			redLock, err := NewRedisMutex(red, "random", 1) //
			if err != nil {
				panic(err)
			}
			if err := redLock.LockWithTimeout(2 * time.Second); err != nil {
				t.Errorf("get lock [%v] failed %v\n", redLock.Key(), err)
				return
			}
			defer redLock.Unlock() // 操作完成后立刻揭开锁
			// defer操作是由下至上出栈被调用
			t.Logf("get lock")
		}()
	}
	wg.Wait()
}
