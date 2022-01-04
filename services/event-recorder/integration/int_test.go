package integration

import (
	"io"
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"ponglehub.co.uk/events/recorder/pkg/recorder"
	"ponglehub.co.uk/lib/events"
)

func assertErr(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Not expecting error: %+v", err)
		t.FailNow()
	}
}

func TestRecorder(t *testing.T) {
	logrus.SetOutput(io.Discard)

	SERVER_URL := os.Getenv("SERVER_URL")

	client, err := events.New(events.EventsArgs{
		BrokerEnv: "BROKER_URL",
		Source:    "int-test",
	})
	assertErr(t, err)

	for _, test := range []struct {
		name   string
		input  []string
		output []string
	}{
		{
			name:   "empty",
			output: []string{},
		},
		{
			name:   "one event",
			input:  []string{"event.test.1"},
			output: []string{"event.test.1"},
		},
		{
			name:   "two events",
			input:  []string{"event.test.1", "event.test.2"},
			output: []string{"event.test.1", "event.test.2"},
		},
	} {
		t.Run(test.name, func(u *testing.T) {
			recorder.Clear(u, SERVER_URL)

			for _, event := range test.input {
				assertErr(u, client.Send(event, "some-data"))
			}
			assert.Equal(u, test.output, recorder.GetEvents(u, SERVER_URL))
		})
	}

	t.Run("latest", func(u *testing.T) {
		recorder.Clear(u, SERVER_URL)
		assertErr(u, client.Send("event.latest", "latest-data"))
		eventType, data := recorder.GetLatest(u, SERVER_URL)
		assert.Equal(u, "event.latest", eventType)
		assert.Equal(u, "\"latest-data\"", data)
	})

	t.Run("wait for latest", func(u *testing.T) {
		recorder.Clear(u, SERVER_URL)

		assertErr(u, client.Send("event.first", "first-data"))

		testOutput := make(chan string, 1)
		go func() {
			testOutput <- recorder.WaitForEvent(u, SERVER_URL, "event.wait")
		}()

		time.Sleep(500 * time.Millisecond)

		assertErr(u, client.Send("event.wait", "wait-data"))

		select {
		case data := <-testOutput:
			assert.Equal(u, data, "\"wait-data\"")
		case <-time.After(5 * time.Second):
			u.Errorf("timed out waiting for event")
			u.Fail()
		}
	})
}
