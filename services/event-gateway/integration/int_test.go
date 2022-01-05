package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"os"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
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

func waitForKey(t *testing.T, rdb *redis.Client, key string) string {
	resultChan := make(chan string, 1)

	go func(resultChan chan<- string) {
		for {
			value, err := rdb.Get(context.Background(), key).Result()
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

type TestClient struct {
	client *http.Client
}

func (t *TestClient) Post(u *testing.T, url string, data map[string]string) []*http.Cookie {
	json_data, err := json.Marshal(data)
	noErr(u, err)

	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(json_data),
	)
	noErr(u, err)

	req.Header.Add("Content-Type", "application/json")

	res, err := t.client.Do(req)
	noErr(u, err)

	assert.Equal(u, 200, res.StatusCode)

	return res.Cookies()
}

func TestGateway(t *testing.T) {
	crds.AddToScheme(scheme.Scheme)

	crdClient, err := crds.New(&crds.ClientArgs{
		External: true,
	})
	noErr(t, err)

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	crdClient.Delete("test-user")
	user, err := crdClient.Create(crds.User{
		Name:    "test-user",
		Display: "test user",
		Email:   "test@user.com",
	})
	noErr(t, err)

	invite := waitForKey(t, rdb, fmt.Sprintf("%s.%s", user.ID, "invite"))

	jar, err := cookiejar.New(nil)
	noErr(t, err)

	testClient := TestClient{
		client: &http.Client{
			Jar: jar,
		},
	}

	testClient.Post(
		t,
		fmt.Sprintf("%s/auth/set-password", os.Getenv("GATEWAY_URL")),
		map[string]string{
			"invite":   invite,
			"password": "new-password",
			"confirm":  "new-password",
		},
	)

	testClient.Post(
		t,
		fmt.Sprintf("%s/auth/login", os.Getenv("GATEWAY_URL")),
		map[string]string{
			"email":    "test@user.com",
			"password": "new-password",
		},
	)

	client, err := events.New(events.EventsArgs{
		BrokerEnv: "GATEWAY_EVENTS",
		Source:    "int-tests",
		Cookies:   jar,
	})
	noErr(t, err)

	err = client.Send("test.event", "event 1")
	noErr(t, err)

	err = client.Send("test.event", "event 2")
	noErr(t, err)

	err = client.Send("test.event", "event 3")
	noErr(t, err)
}
