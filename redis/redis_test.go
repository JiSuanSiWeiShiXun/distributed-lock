package redis

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-redis/redis/v8"
)

var (
	red *redis.Client
)

func setUp() {
	red = redis.NewClient(&redis.Options{
		Addr:     "124.222.247.58:6379",
		Password: "",
	})
	fmt.Printf("redis ping return [%v]\n", red.Ping(context.Background()).Val())
}

func tearDown() {

}

// TestRunScript golang调用lua脚本实现redis事务的逻辑
func TestRunScript(t *testing.T) {
	setUp()
	res, err := tmpTransaction.Run(context.Background(), red, []string{"key1", "key2"}, "111", "222").Result()
	if res.(string) != "OK" {
		t.Fatalf("[err] %v [res] [%s]\n", err, res.(string))
	}
	t.Logf("[res] %v\n", res)
}
