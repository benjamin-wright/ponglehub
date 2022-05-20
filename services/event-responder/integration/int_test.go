package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"ponglehub.co.uk/lib/events"
)

func noErr(t *testing.T, err error) {
	if err != nil {
		assert.NoError(t, err)
		t.FailNow()
	}
}

func clearEvents(t *testing.T, rdb *redis.Client, id string) {
	noErr(t, rdb.Del(context.Background(), fmt.Sprintf("%s.responses", id)).Err())
}

func pubsubChannel(t *testing.T, rdb *redis.Client, id string) <-chan *redis.Message {
	pubsub := rdb.Subscribe(context.TODO(), fmt.Sprintf("%s.responses", id))
	return pubsub.Channel(redis.WithChannelSize(10))
}

func TestEvents(t *testing.T) {
	client, err := events.New(events.EventsArgs{
		BrokerEnv: "RESPONDER_URL",
		Source:    "test",
	})
	noErr(t, err)

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	type testEvent struct {
		Type   string
		Data   interface{}
		UserId interface{}
	}

	const TEST_USER = "1234"
	const OTHER_USER = "5678"

	for _, test := range []struct {
		name     string
		events   []testEvent
		expected []map[string]interface{}
	}{
		{
			name: "single",
			events: []testEvent{
				{Type: "test.event", Data: "messages", UserId: TEST_USER},
			},
			expected: []map[string]interface{}{
				{"data": "\"messages\"", "type": "test.event"},
			},
		},
		{
			name: "double",
			events: []testEvent{
				{Type: "test.event", Data: "message 1", UserId: TEST_USER},
				{Type: "another.event", Data: "message 2", UserId: TEST_USER},
			},
			expected: []map[string]interface{}{
				{"data": "\"message 1\"", "type": "test.event"},
				{"data": "\"message 2\"", "type": "another.event"},
			},
		},
		{
			name: "one from another user",
			events: []testEvent{
				{Type: "test.event", Data: "message 1", UserId: TEST_USER},
				{Type: "another.event", Data: "message 2", UserId: OTHER_USER},
			},
			expected: []map[string]interface{}{
				{"data": "\"message 1\"", "type": "test.event"},
			},
		},
		{
			name: "other one from another user",
			events: []testEvent{
				{Type: "test.event", Data: "message 1", UserId: OTHER_USER},
				{Type: "another.event", Data: "message 2", UserId: TEST_USER},
			},
			expected: []map[string]interface{}{
				{"data": "\"message 2\"", "type": "another.event"},
			},
		},
	} {
		t.Run(test.name, func(u *testing.T) {
			clearEvents(u, rdb, TEST_USER)
			clearEvents(u, rdb, OTHER_USER)

			pubsub := rdb.Subscribe(context.TODO(), TEST_USER+".responses")
			defer pubsub.Close()
			responseChannel := pubsub.Channel(redis.WithChannelSize(10))

			time.Sleep(time.Millisecond * 500)

			for _, event := range test.events {
				client.Send(
					event.Type,
					event.Data,
					map[string]interface{}{"userid": event.UserId},
				)
			}

			for _, expected := range test.expected {
				select {
				case actual := <-responseChannel:
					data, _ := json.Marshal(expected)
					assert.Equal(u, string(data), actual.Payload)
				case <-time.After(time.Second * 2):
					assert.FailNow(u, "timed out waiting for event")
				}
			}

			select {
			case <-responseChannel:
				assert.FailNow(u, "received extra event")
			case <-time.After(time.Second):
			}
		})
	}
}
