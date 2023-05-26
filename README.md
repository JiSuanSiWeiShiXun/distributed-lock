# distributed-lock

## redis实现的分布式锁
目前只支持单redis实例的分布式锁，不支持红锁要求的多redis实例的实现原则
提供RedisMutex类，在其对象上实现了以下接口
1. Lock 获取锁, 失败立即返回
2. LockWithTimeout: 获取锁, 以间隔100ms频率获取锁，直到超时或者获取成功
3. Unlock: 释放锁

## example
```Golang
import (
    redlock "github.com/JiSuanSiWeiShiXun/distributed-lock/redis"
)

// Transaction 需要加锁的事务操作
func Transaction() {
    // 加锁
	redMutex, err := redlock.NewRedisMutex(db.Red(), taskID, 1)
	if err != nil {
		logs.Logger.Errorf("get redmutex of [key] %v failed: %v", redMutex.Key(), err)
		c.JSON(http.StatusInternalServerError, "server busy, try it again later")
		return
	}
	redMutex.LockWithTimeout(3 * time.Second)
	defer redMutex.Unlock()
    
    ... ...
}
```

## TODO
1. 锁续期，到一半超时时间的时候续期
2. 多redis实例
3. 其他分布式锁实现方式