package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"ponglehub.co.uk/auth/auth-server/internal/client"
	"ponglehub.co.uk/auth/auth-server/internal/testutils"
)

func TestServer(t *testing.T) {
	logrus.SetOutput(io.Discard)

	testUrl, ok := os.LookupEnv("TEST_SERVER")
	if !ok {
		t.Error("TEST_SERVER env var not found")
		t.FailNow()
	}

	cli := testutils.Client(t)
	defer cli.Close(context.TODO())

	for _, test := range []struct {
		name     string
		email    string
		password string
		code     int
	}{
		{name: "success", email: "johnny@place.com", password: "some-pass", code: http.StatusOK},
		{name: "no user", email: "wrong@place.com", password: "some-pass", code: http.StatusUnauthorized},
		{name: "wrong password", email: "johhny@place.com", password: "wrong-pass", code: http.StatusUnauthorized},
		{name: "no password", email: "johhny@place.com", code: http.StatusBadRequest},
		{name: "no user", password: "some-pass", code: http.StatusBadRequest},
		{name: "empty request", code: http.StatusBadRequest},
	} {
		t.Run(test.name, func(u *testing.T) {
			if err := cli.Clear(); err != nil {
				u.Fatalf("Error clearing database: %+v", err)
			}

			_, err := cli.AddUser(context.TODO(), client.User{
				Name:     "johnny",
				Email:    "johnny@place.com",
				Password: testutils.Hash(u, "some-pass"),
				Verified: true,
			})
			if err != nil {
				fmt.Printf("Failed to add user: %+v\n", err)
				u.Fail()
				return
			}

			data, _ := json.Marshal(map[string]string{
				"email":    test.email,
				"password": test.password,
			})
			resp, err := http.Post(testUrl+"/", "application/json", bytes.NewBuffer(data))
			if err != nil {
				assert.FailNow(u, "Failed making test request: %+v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != test.code {
				fmt.Printf("Expected %d: Recieved %d\n", test.code, resp.StatusCode)
				u.Fail()
			}
		})
	}
}
