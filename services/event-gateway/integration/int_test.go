package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/scheme"
	"ponglehub.co.uk/events/gateway/internal/crds"
	"ponglehub.co.uk/lib/events"
)

func noErr(t *testing.T, err error) {
	if err != nil {
		assert.NoError(t, err)
		t.FailNow()
	}
}

func TestGateway(t *testing.T) {
	crds.AddToScheme(scheme.Scheme)

	crdClient, err := crds.New(&crds.ClientArgs{
		External: true,
	})
	noErr(t, err)

	client, err := events.New(events.EventsArgs{
		BrokerEnv: "GATEWAY_URL",
		Source:    "int-tests",
	})
	noErr(t, err)

	crdClient.Delete("test-user")
	_, err = crdClient.Create(crds.User{
		Name:    "test-user",
		Display: "test user",
		Email:   "test@user.com",
	})
	noErr(t, err)

	noErr(t, client.Send("test.event", "some data"))
}
