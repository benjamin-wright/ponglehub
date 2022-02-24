package redis

import (
	"context"
	"fmt"
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

func (r *Redis) AddResponses(t *testing.T, id string, responses []interface{}) {
	if len(responses) == 0 {
		return
	}

	key := fmt.Sprintf("%s.responses", id)
	err := r.client.RPush(context.Background(), key, responses...).Err()
	if err != nil {
		assert.FailNow(t, "failed to add responses", err)
	}
}

func (r *Redis) ClearResponses(t *testing.T, id string) {
	key := fmt.Sprintf("%s.responses", id)
	err := r.client.Del(context.Background(), key).Err()
	if err != nil {
		assert.FailNow(t, "failed to clear responses", err)
	}
}

func (r *Redis) GetResponses(t *testing.T, id string) []string {
	key := fmt.Sprintf("%s.responses", id)
	values, err := r.client.LRange(context.Background(), key, 0, -1).Result()
	if err == redis.Nil {
		return nil
	} else if err != nil {
		assert.Fail(t, "failed to fetch responses:", err)
	}

	return values
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
