package integration

import (
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/scheme"
	"ponglehub.co.uk/events/gateway/integration/redis"
	"ponglehub.co.uk/events/gateway/integration/test_client"
	"ponglehub.co.uk/events/gateway/internal/services/crds"
	"ponglehub.co.uk/events/recorder/pkg/recorder"
	"ponglehub.co.uk/lib/events"
)

func init() {
	logrus.SetOutput(io.Discard)
	crds.AddToScheme(scheme.Scheme)
}

func noErr(t *testing.T, err error) {
	if err != nil {
		assert.NoError(t, err)
		t.FailNow()
	}
}

func clients(t *testing.T) (*crds.UserClient, *redis.Redis, *test_client.TestClient, *events.Events) {
	crdClient, err := crds.New(&crds.ClientArgs{
		External: true,
	})
	noErr(t, err)

	redis := redis.New()

	testClient := test_client.New(t)

	eventClient, err := events.New(events.EventsArgs{
		BrokerEnv: "GATEWAY_EVENTS",
		Source:    "int-tests",
		Cookies:   testClient.CookieJar(),
	})
	noErr(t, err)

	return crdClient, redis, testClient, eventClient
}

func TestInviteToken(t *testing.T) {
	for _, test := range []struct {
		Name       string
		Prep       func(*testing.T, *redis.Redis, crds.User, *test_client.TestClient)
		Input      func(string) map[string]string
		StatusCode int
	}{
		{
			Name: "success",
			Input: func(invite string) map[string]string {
				return map[string]string{
					"invite":   invite,
					"password": "new-password",
					"confirm":  "new-password",
				}
			},
			StatusCode: 200,
		},
		{
			Name: "expired",
			Prep: func(t *testing.T, r *redis.Redis, u crds.User, c *test_client.TestClient) {
				r.DeleteKey(t, fmt.Sprintf("%s.%s", u.ID, "invite"))
			},
			Input: func(invite string) map[string]string {
				return map[string]string{
					"invite":   invite,
					"password": "new-password",
					"confirm":  "new-password",
				}
			},
			StatusCode: 401,
		},
		{
			Name: "expired",
			Prep: func(t *testing.T, r *redis.Redis, u crds.User, c *test_client.TestClient) {
				invite := r.WaitForKey(t, fmt.Sprintf("%s.%s", u.ID, "invite"))

				url := fmt.Sprintf("%s/auth/set-password", os.Getenv("GATEWAY_URL"))
				res := c.Post(
					t,
					url,
					map[string]string{
						"invite":   invite,
						"password": "new-password",
						"confirm":  "new-password",
					},
				)
				assert.Equal(t, 200, res.StatusCode)
			},
			Input: func(invite string) map[string]string {
				return map[string]string{
					"invite":   invite,
					"password": "new-password",
					"confirm":  "new-password",
				}
			},
			StatusCode: 401,
		},
		{
			Name:       "no args",
			Input:      func(invite string) map[string]string { return map[string]string{} },
			StatusCode: 400,
		},
		{
			Name: "malformed token",
			Input: func(invite string) map[string]string {
				return map[string]string{
					"invite":   "bad-token",
					"password": "new-password",
					"confirm":  "new-password",
				}
			},
			StatusCode: 401,
		},
		{
			Name: "mismatched passwords",
			Input: func(invite string) map[string]string {
				return map[string]string{
					"invite":   "bad-token",
					"password": "new-password",
					"confirm":  "wrong-password",
				}
			},
			StatusCode: 400,
		},
	} {
		t.Run(test.Name, func(u *testing.T) {
			crdClient, redisClient, testClient, _ := clients(u)

			crdClient.Delete("test-user")
			user, err := crdClient.Create(crds.User{
				Name:    "test-user",
				Display: "test user",
				Email:   "test@user.com",
			})
			noErr(u, err)

			invite := redisClient.WaitForKey(u, fmt.Sprintf("%s.%s", user.ID, "invite"))

			if test.Prep != nil {
				test.Prep(u, redisClient, user, testClient)
			}

			url := fmt.Sprintf("%s/auth/set-password", os.Getenv("GATEWAY_URL"))
			res := testClient.Post(u, url, test.Input(invite))

			assert.Equal(u, test.StatusCode, res.StatusCode)
		})
	}
}

func TestLogin(t *testing.T) {
	for _, test := range []struct {
		Name       string
		Input      map[string]string
		StatusCode int
		Cookies    int
	}{
		{
			Name: "success",
			Input: map[string]string{
				"email":    "test@user.com",
				"password": "new-password",
			},
			StatusCode: 200,
			Cookies:    1,
		},
		{
			Name:       "no args",
			Input:      map[string]string{},
			StatusCode: 400,
			Cookies:    0,
		},
		{
			Name: "wrong email",
			Input: map[string]string{
				"email":    "wrong@user.com",
				"password": "new-password",
			},
			StatusCode: 401,
			Cookies:    0,
		},
		{
			Name: "wrong password",
			Input: map[string]string{
				"email":    "test@user.com",
				"password": "wrong-password",
			},
			StatusCode: 401,
			Cookies:    0,
		},
	} {
		t.Run(test.Name, func(u *testing.T) {
			crdClient, redisClient, testClient, _ := clients(u)

			crdClient.Delete("test-user")
			user, err := crdClient.Create(crds.User{
				Name:    "test-user",
				Display: "test user",
				Email:   "test@user.com",
			})
			noErr(u, err)

			invite := redisClient.WaitForKey(u, fmt.Sprintf("%s.%s", user.ID, "invite"))

			url := fmt.Sprintf("%s/auth/set-password", os.Getenv("GATEWAY_URL"))
			res := testClient.Post(
				u,
				url,
				map[string]string{
					"invite":   invite,
					"password": "new-password",
					"confirm":  "new-password",
				},
			)
			assert.Equal(u, 200, res.StatusCode)

			url = fmt.Sprintf("%s/auth/login", os.Getenv("GATEWAY_URL"))
			res = testClient.Post(t, url, test.Input)
			assert.Equal(t, test.StatusCode, res.StatusCode)
			assert.Equal(t, test.Cookies, len(res.Cookies()))
		})
	}
}

func TestProxying(t *testing.T) {
	for _, test := range []struct {
		Name         string
		LoggedIn     bool
		Unauthorized bool
		Events       int
	}{
		{
			Name:         "not logged in",
			LoggedIn:     false,
			Unauthorized: true,
			Events:       0,
		},
		{
			Name:         "logged in",
			LoggedIn:     true,
			Unauthorized: false,
			Events:       1,
		},
	} {
		t.Run(test.Name, func(u *testing.T) {
			crdClient, redisClient, testClient, eventClient := clients(u)
			recorder.Clear(u, os.Getenv("RECORDER_URL"))

			crdClient.Delete("test-user")
			user, err := crdClient.Create(crds.User{
				Name:    "test-user",
				Display: "test user",
				Email:   "test@user.com",
			})
			noErr(u, err)

			invite := redisClient.WaitForKey(u, fmt.Sprintf("%s.%s", user.ID, "invite"))

			url := fmt.Sprintf("%s/auth/set-password", os.Getenv("GATEWAY_URL"))
			res := testClient.Post(
				u,
				url,
				map[string]string{
					"invite":   invite,
					"password": "new-password",
					"confirm":  "new-password",
				},
			)
			assert.Equal(u, 200, res.StatusCode)

			if test.LoggedIn {
				url = fmt.Sprintf("%s/auth/login", os.Getenv("GATEWAY_URL"))
				res = testClient.Post(
					t,
					url,
					map[string]string{
						"email":    "test@user.com",
						"password": "new-password",
					},
				)
				assert.Equal(t, 200, res.StatusCode)
			}

			err = eventClient.Send("test.event", "event 1")
			if test.Unauthorized {
				assert.Equal(t, events.UnauthorizedError, err)
			} else {
				noErr(t, err)
			}

			time.Sleep(250 * time.Millisecond)
			received := recorder.GetEvents(u, os.Getenv("RECORDER_URL"))
			assert.Equal(t, test.Events, len(received))
		})
	}
}
