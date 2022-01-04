package integration

import (
	"os"
	"testing"

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

	client.Delete("test-trigger", TEST_NAMESPACE)
	err = client.Create("test-trigger", TEST_NAMESPACE, []string{"test.event"}, RECORDER_URL)
	assert.NoError(t, err)

	recorder.Clear(t, RECORDER_SERVER)

	sender, err := events.New(events.EventsArgs{
		BrokerURL: BROKER_URL,
		Source:    "unit-tests",
	})
	if err != nil {
		assert.NoError(t, err)
		t.FailNow()
	}

	err = sender.Send("test.event", "some event data")
	if err != nil {
		assert.NoError(t, err)
		t.FailNow()
	}

	t.Logf("Waiting for event...")
	events := recorder.WaitForEvents(t, RECORDER_SERVER, 1)
	assert.Equal(t, []string{"test.event"}, events)
}
