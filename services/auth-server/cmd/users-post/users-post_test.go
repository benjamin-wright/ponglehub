package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/yaml.v2"
	main "ponglehub.co.uk/auth/auth-server/cmd/users-post"
	"ponglehub.co.uk/auth/auth-server/internal/client"
	"ponglehub.co.uk/auth/auth-server/internal/server"
	"ponglehub.co.uk/auth/auth-server/internal/testutils"
)

const TEST_DB = "test_users_post"

func loadUsers(t *testing.T, cli *testutils.TestClient, exists bool) {
	if err := cli.Reset(); err != nil {
		t.Fatalf("Error clearing database: %+v", err)
	}

	if exists {
		if err := cli.AddUser("123e4567-e89b-12d3-a456-426614174000", "username", "user@email.com", "some-pwd", false); err != nil {
			t.Fatalf("Error adding test user: %+v", err)
		}
	}
}

func TestPostRoute(t *testing.T) {
	cli, err := testutils.NewClient(TEST_DB)
	if err != nil {
		fmt.Printf("Failed to create test client: %+v\n", err)
		t.Fail()
		return
	}
	defer cli.Close()

	for _, test := range []struct {
		name     string
		args     map[string]string
		code     int
		exists   bool
		expected *client.User
	}{
		{
			name:     "correct",
			args:     map[string]string{"name": "username", "email": "user@email.com", "password": "user-password"},
			code:     http.StatusAccepted,
			expected: &client.User{Name: "username", Email: "user@email.com", Password: "user-password", Verified: false},
		},
		{
			name: "missing name",
			args: map[string]string{"email": "user@email.com", "password": "user-password"},
			code: http.StatusBadRequest,
		},
		{
			name: "missing email",
			args: map[string]string{"name": "username", "password": "user-password"},
			code: http.StatusBadRequest,
		},
		{
			name: "missing password",
			args: map[string]string{"name": "username", "email": "user@email.com"},
			code: http.StatusBadRequest,
		},
		{
			name:   "already exists",
			args:   map[string]string{"name": "username", "email": "user@email.com", "password": "user-password"},
			code:   http.StatusConflict,
			exists: true,
		},
		{
			name:   "same username",
			args:   map[string]string{"name": "username", "email": "different@email.com", "password": "user-password"},
			code:   http.StatusConflict,
			exists: true,
		},
		{
			name:   "same email",
			args:   map[string]string{"name": "different", "email": "user@email.com", "password": "user-password"},
			code:   http.StatusConflict,
			exists: true,
		},
		{
			name:     "different user",
			args:     map[string]string{"name": "different", "email": "different@email.com", "password": "user-password"},
			code:     http.StatusAccepted,
			exists:   true,
			expected: &client.User{Name: "different", Email: "different@email.com", Password: "user-password", Verified: false},
		},
	} {
		t.Run(test.name, func(u *testing.T) {
			loadUsers(u, cli, test.exists)

			r := server.GetRouter(cli.TargetConfig(), main.RouteBuilder)

			data, _ := json.Marshal(test.args)

			req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(data))
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != test.code {
				fmt.Printf("Expected %d: Recieved %d\n", test.code, w.Code)
				u.Fail()
			}

			if test.expected != nil {
				var result struct {
					ID string `json:"id"`
				}

				err = yaml.Unmarshal(w.Body.Bytes(), &result)
				if err != nil {
					u.Fatalf("Error unmarshalling post response: %v", err)
				}

				user, err := cli.GetUser(result.ID)
				if err != nil {
					u.Fatalf("Error fetching new user by id: %v", err)
				}

				if user.Name != test.expected.Name || user.Email != test.expected.Email || user.Verified != test.expected.Verified {
					fmt.Printf("Expected %+v to equal %+v", user, test.expected)
					u.Fail()
				}

				err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(test.expected.Password))
				if err != nil {
					fmt.Printf("User failed password check: %+v", err)
					u.Fail()
				}
			}
		})
	}
}
