package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/crypto/bcrypt"
	main "ponglehub.co.uk/auth/auth-server/cmd/users-put"
	"ponglehub.co.uk/auth/auth-server/internal/client"
	"ponglehub.co.uk/auth/auth-server/internal/server"
	"ponglehub.co.uk/auth/auth-server/internal/testutils"
)

const TEST_DB = "test_users_put"

func loadUsers(t *testing.T, cli *testutils.TestClient, existing *client.User) {
	if err := cli.Reset(); err != nil {
		t.Fatalf("Error clearing database: %+v", err)
	}

	if existing != nil {
		if err := cli.AddUser(existing.ID, existing.Name, existing.Email, existing.Password, existing.Verified); err != nil {
			t.Fatalf("Error adding test user: %+v", err)
		}
	}
}

func TestPutRoute(t *testing.T) {
	cli, err := testutils.NewClient(TEST_DB)
	if err != nil {
		fmt.Printf("Failed to create test client: %+v\n", err)
		t.Fail()
		return
	}
	defer cli.Close()

	oldUser := client.User{
		ID:       "123e4567-e89b-12d3-a456-426614174001",
		Name:     "old-name",
		Email:    "old-email",
		Password: "old-password",
		Verified: false,
	}

	for _, test := range []struct {
		name     string
		path     string
		args     map[string]interface{}
		code     int
		existing *client.User
		expected *client.User
	}{
		{
			name:     "bad id",
			path:     "/whatevs",
			args:     map[string]interface{}{"name": "new-name"},
			code:     http.StatusBadRequest,
			existing: &oldUser,
		},
		{
			name:     "missing args",
			path:     "/123e4567-e89b-12d3-a456-426614174001",
			args:     map[string]interface{}{},
			code:     http.StatusBadRequest,
			existing: &oldUser,
		},
		{
			name:     "wrong id",
			path:     "/123e4567-e89b-12d3-a456-426614174000",
			args:     map[string]interface{}{"name": "new-name"},
			code:     http.StatusNotFound,
			existing: &oldUser,
		},
		{
			name:     "update name",
			path:     "/123e4567-e89b-12d3-a456-426614174001",
			args:     map[string]interface{}{"name": "new-name"},
			code:     http.StatusAccepted,
			existing: &oldUser,
			expected: &client.User{ID: "123e4567-e89b-12d3-a456-426614174001", Name: "new-name", Email: "old-email", Password: "old-password", Verified: false},
		},
		{
			name:     "update email",
			path:     "/123e4567-e89b-12d3-a456-426614174001",
			args:     map[string]interface{}{"email": "new-email"},
			code:     http.StatusAccepted,
			existing: &oldUser,
			expected: &client.User{ID: "123e4567-e89b-12d3-a456-426614174001", Name: "old-name", Email: "new-email", Password: "old-password", Verified: false},
		},
		{
			name:     "update password",
			path:     "/123e4567-e89b-12d3-a456-426614174001",
			args:     map[string]interface{}{"password": "new-password"},
			code:     http.StatusAccepted,
			existing: &oldUser,
			expected: &client.User{ID: "123e4567-e89b-12d3-a456-426614174001", Name: "old-name", Email: "old-email", Password: "new-password", Verified: false},
		},
		{
			name:     "update verified",
			path:     "/123e4567-e89b-12d3-a456-426614174001",
			args:     map[string]interface{}{"verified": true},
			code:     http.StatusAccepted,
			existing: &oldUser,
			expected: &client.User{ID: "123e4567-e89b-12d3-a456-426614174001", Name: "old-name", Email: "old-email", Password: "old-password", Verified: true},
		},
		{
			name:     "don't disable verified if not provided",
			path:     "/123e4567-e89b-12d3-a456-426614174001",
			args:     map[string]interface{}{"name": "new-name"},
			code:     http.StatusAccepted,
			existing: &client.User{ID: "123e4567-e89b-12d3-a456-426614174001", Name: "old-name", Email: "old-email", Password: "old-password", Verified: true},
			expected: &client.User{ID: "123e4567-e89b-12d3-a456-426614174001", Name: "new-name", Email: "old-email", Password: "old-password", Verified: true},
		},
		{
			name:     "everything at once",
			path:     "/123e4567-e89b-12d3-a456-426614174001",
			args:     map[string]interface{}{"name": "new-name", "email": "new-email", "password": "new-password", "verified": true},
			code:     http.StatusAccepted,
			existing: &oldUser,
			expected: &client.User{ID: "123e4567-e89b-12d3-a456-426614174001", Name: "new-name", Email: "new-email", Password: "new-password", Verified: true},
		},
	} {
		t.Run(test.name, func(u *testing.T) {
			loadUsers(u, cli, test.existing)

			r := server.GetRouter(cli.TargetConfig(), main.RouteBuilder)

			data, _ := json.Marshal(test.args)

			req, _ := http.NewRequest("PUT", test.path, bytes.NewBuffer(data))
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != test.code {
				fmt.Printf("Expected %d: Recieved %d\n", test.code, w.Code)
				u.Fail()
			}

			if test.expected != nil {
				user, err := cli.GetUser(test.expected.ID)
				if err != nil {
					u.Fatalf("Failed to get user: %+v", err)
				}

				if user.Name != test.expected.Name || user.Email != test.expected.Email || user.Verified != test.expected.Verified {
					fmt.Printf("Expected %+v to equal %+v", user, *test.expected)
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
