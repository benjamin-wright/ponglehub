package tests

import (
	"io"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/scheme"
	"ponglehub.co.uk/events/broker/internal/crds"
	"ponglehub.co.uk/events/recorder/pkg/recorder"
	"ponglehub.co.uk/lib/events"
)

func getEnv(t *testing.T, env string) string {
	value, ok := os.LookupEnv(env)
	if !ok {
		t.Logf("Missing environment variable %s", env)
		t.FailNow()
	}

	return value
}

func TestSomething(t *testing.T) {
	logrus.SetOutput(io.Discard)

	TEST_NAMESPACE := getEnv(t, "TEST_NAMESPACE")
	BROKER_URL := getEnv(t, "BROKER_URL")
	RECORDER_URL := getEnv(t, "RECORDER_URL")
	RECORDER_SERVER := getEnv(t, "RECORDER_SERVER")

	crds.AddToScheme(scheme.Scheme)
	client, err := crds.New(&crds.ClientArgs{External: true})
	if err != nil {
		assert.NoError(t, err)
		t.FailNow()
	}

	sender, err := events.New(events.EventsArgs{
		BrokerURL: BROKER_URL,
		Source:    "unit-tests",
	})
	if err != nil {
		assert.NoError(t, err)
		t.FailNow()
	}

	for _, test := range []struct {
		name     string
		filters  []string
		expected []string
	}{
		{
			name:     "normal filtering",
			filters:  []string{"test.event"},
			expected: []string{"test.event"},
		},
		{
			name:     "mutliple filters",
			filters:  []string{"test.event", "other.event"},
			expected: []string{"test.event", "other.event"},
		},
		{
			name:     "closed wildcard filtering",
			filters:  []string{"test.*.resp"},
			expected: []string{"test.event.resp", "test.random.resp"},
		},
		{
			name:     "open wildcard filtering",
			filters:  []string{"test.*"},
			expected: []string{"test.event", "test.event.resp", "test.event.random", "test.random", "test.random.resp"},
		},
		{
			name:     "no filters",
			filters:  []string{},
			expected: []string{},
		},
	} {
		t.Run(test.name, func(u *testing.T) {
			client.Delete("test-trigger", TEST_NAMESPACE)
			err = client.Create("test-trigger", TEST_NAMESPACE, test.filters, RECORDER_URL)
			assert.NoError(u, err)

			recorder.Clear(u, RECORDER_SERVER)

			for _, event := range []string{"test", "test.event", "test.event.resp", "test.event.random", "test.random", "test.random.resp", "other", "other.event", "other.event.resp"} {
				err = sender.Send(event, "some event data")
				if err != nil {
					assert.NoError(u, err)
					u.FailNow()
				}
			}

			events := recorder.WaitForEvents(u, RECORDER_SERVER, len(test.expected))
			assert.Equal(u, test.expected, events)
		})
	}
}
