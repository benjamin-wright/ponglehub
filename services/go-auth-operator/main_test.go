package main_test

import (
	"context"
	"os"
	"testing"

	cenats "github.com/cloudevents/sdk-go/protocol/nats/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"ponglehub.co.uk/auth/auth-operator/internal/client"
)

func getEnv(t *testing.T, key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		t.Errorf("Missing environment variable %s", key)
		t.FailNow()
	}

	return value
}

func TestAddCRD(t *testing.T) {
	client.AddToScheme(scheme.Scheme)

	natsUrl := getEnv(t, "NATS_URL")
	natsSubject := getEnv(t, "NATS_SUBJECT")

	for _, test := range []struct {
		Name string
	}{
		{
			Name: "success",
		},
	} {
		t.Run(test.Name, func(u *testing.T) {
			p, err := cenats.NewConsumer(natsUrl, natsSubject, cenats.NatsOptions())
			if err != nil {
				assert.FailNow(u, "Error creating test receiver: %+v", err)
			}

			events, err := cloudevents.NewClient(p)
			if err != nil {
				assert.FailNow(u, "error creating cloudevents instance: %+v", err)
			}

			received := make(chan event.Event)

			ctx, cancel := context.WithCancel(context.Background())
			go func() {
				err = events.StartReceiver(ctx, func(e event.Event) {
					received <- e
				})

				if err != nil {
					assert.FailNow(u, "Failed to start received: %+v", err)
				}
			}()
			defer cancel()

			cli, err := client.New()
			if err != nil {
				assert.FailNow(u, "Failed to start client: %+v", err)
			}

			err = cli.Create(client.AuthUser{
				ObjectMeta: v1.ObjectMeta{
					Name: "test-user",
				},
				Spec: client.AuthUserSpec{
					Name:  "test-user",
					Email: "test@user.com",
				},
			}, v1.CreateOptions{})

			if err != nil {
				assert.FailNow(u, "failed to create user: %+v", err)
			}

			e := <-received

			assert.Equal(u, "", e.String())
		})
	}
}
