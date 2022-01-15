package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	client *redis.Client
}

func New(url string) *Redis {
	rdb := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &Redis{client: rdb}
}

func (r *Redis) DeleteKey(key string) error {
	err := r.client.Del(context.Background(), key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete redis key: %+v", err)
	}

	return nil
}

func (r *Redis) ListKeys() ([]string, error) {
	return r.client.Keys(context.Background(), "*").Result()
}

func (r *Redis) WaitForKey(key string) (string, error) {
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
		return result, nil
	case <-time.After(5 * time.Second):
		return "", fmt.Errorf("timed out waiting for key %s", key)
	}
}
