package integration

import (
	"os"
	"testing"

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

}
