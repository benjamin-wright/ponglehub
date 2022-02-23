package storage

import (
	"context"
	"fmt"

	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/go-redis/redis/v8"
)

type Storage struct {
	redis *redis.Client
}

func New(redisUrl string) (*Storage, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisUrl,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	storage := Storage{
		redis: rdb,
	}

	return &storage, nil
}

func (s *Storage) AddEvent(id string, event event.Event) error {
	key := fmt.Sprintf("%s.responses", id)
	data, err := event.MarshalJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal event to json: %+v", err)
	}

	err = s.redis.SAdd(context.Background(), key, string(data)).Err()
	if err != nil {
		return fmt.Errorf("failed to fetch add message to redis store: %+v", err)
	}

	return nil
}
