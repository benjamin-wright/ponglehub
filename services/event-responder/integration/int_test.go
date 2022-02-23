package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/cloudevents/sdk-go/v2/event"
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

func getEvents(t *testing.T, rdb *redis.Client, id string) []string {
	members, err := rdb.LRange(context.Background(), fmt.Sprintf("%s.responses", id), 0, -1).Result()
	noErr(t, err)

	return members
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
		expected []testEvent
	}{
		{
			name: "single",
			events: []testEvent{
				{Type: "test.event", Data: "messages", UserId: TEST_USER},
			},
			expected: []testEvent{
				{Type: "test.event", Data: "messages", UserId: TEST_USER},
			},
		},
		{
			name: "double",
			events: []testEvent{
				{Type: "test.event", Data: "message 1", UserId: TEST_USER},
				{Type: "another.event", Data: "message 2", UserId: TEST_USER},
			},
			expected: []testEvent{
				{Type: "test.event", Data: "message 1", UserId: TEST_USER},
				{Type: "another.event", Data: "message 2", UserId: TEST_USER},
			},
		},
		{
			name: "one from another user",
			events: []testEvent{
				{Type: "test.event", Data: "message 1", UserId: TEST_USER},
				{Type: "another.event", Data: "message 2", UserId: OTHER_USER},
			},
			expected: []testEvent{
				{Type: "test.event", Data: "message 1", UserId: TEST_USER},
			},
		},
		{
			name: "other one from another user",
			events: []testEvent{
				{Type: "test.event", Data: "message 1", UserId: OTHER_USER},
				{Type: "another.event", Data: "message 2", UserId: TEST_USER},
			},
			expected: []testEvent{
				{Type: "another.event", Data: "message 2", UserId: TEST_USER},
			},
		},
	} {
		t.Run(test.name, func(u *testing.T) {
			clearEvents(u, rdb, TEST_USER)
			clearEvents(u, rdb, OTHER_USER)

			clearedEvents := getEvents(u, rdb, TEST_USER)
			if !assert.Equal(u, []string{}, clearedEvents) {
				u.FailNow()
			}

			for _, event := range test.events {
				client.Send(
					event.Type,
					event.Data,
					map[string]interface{}{"userid": event.UserId},
				)
			}

			testEvents := getEvents(u, rdb, TEST_USER)
			if !assert.Equal(u, len(test.expected), len(testEvents)) {
				u.FailNow()
			}

			for index, data := range testEvents {
				event := event.Event{}
				noErr(u, event.UnmarshalJSON([]byte(data)))

				var output interface{}
				noErr(u, json.Unmarshal(event.Data(), &output))

				userId, err := event.Context.GetExtension("userid")
				noErr(u, err)

				assert.Equal(u, test.expected[index].Type, event.Type())
				assert.Equal(u, test.expected[index].Data, output)
				assert.Equal(u, test.expected[index].UserId, userId)
			}

		})
	}
}
