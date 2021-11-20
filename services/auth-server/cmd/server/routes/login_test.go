package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"ponglehub.co.uk/auth/auth-server/internal/client"
	"ponglehub.co.uk/auth/auth-server/internal/testutils"
)

const TEST_DB = "test_login"

func TestLoginRoute(t *testing.T) {
	logrus.SetOutput(io.Discard)
	gin.SetMode(gin.TestMode)

	cli := testutils.Client(t, TEST_DB)
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

			w := httptest.NewRecorder()
			ctx, r := gin.CreateTestContext(w)
			r.POST("/", LoginHandler(cli))

			data, _ := json.Marshal(map[string]string{
				"email":    test.email,
				"password": test.password,
			})
			ctx.Request, err = http.NewRequest("POST", "/", bytes.NewBuffer(data))
			if err != nil {
				assert.FailNow(u, "Failed making test request: %+v", err)
			}

			r.HandleContext(ctx)

			if w.Code != test.code {
				fmt.Printf("Expected %d: Recieved %d\n", test.code, w.Code)
				u.Fail()
			}
		})
	}
}
