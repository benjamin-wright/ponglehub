package redis

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

type Redis struct {
	client *redis.Client
}

func New() *Redis {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &Redis{client: rdb}
}

func (r *Redis) DeleteKey(t *testing.T, key string) {
	err := r.client.Del(context.Background(), key).Err()
	if err != nil {
		assert.FailNow(t, "failed to delete key", err)
	}
}

func (r *Redis) WaitForKey(t *testing.T, key string) string {
	resultChan := make(chan string, 1)

	go func(resultChan chan<- string) {
		for {
			value, err := r.client.Get(context.Background(), key).Result()
			if err != nil {
				continue
			}

			resultChan <- value
			break
		}
	}(resultChan)

	select {
	case result := <-resultChan:
		return result
	case <-time.After(5 * time.Second):
		t.Errorf("timed out waiting for key: %s", key)
		t.FailNow()
		return ""
	}
}
