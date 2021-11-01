package main_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	main "ponglehub.co.uk/auth/auth-server/cmd/users-get"
	"ponglehub.co.uk/auth/auth-server/internal/server"
	"ponglehub.co.uk/auth/auth-server/internal/testutils"
)

const TEST_DB = "test_users_get"

func loadUsers(t *testing.T, cli *testutils.TestClient, ids []string) {
	if err := cli.Reset(); err != nil {
		t.Fatalf("Error clearing database: %+v", err)
	}

	for idx, id := range ids {
		if err := cli.AddUser(id, fmt.Sprintf("user-%d", idx), fmt.Sprintf("email-%d", idx), "some-pwd", true); err != nil {
			t.Fatalf("Error adding test user: %+v", err)
		}
	}
}

func TestGetRoute(t *testing.T) {
	cli, err := testutils.NewClient(TEST_DB)
	if err != nil {
		fmt.Printf("Failed to create test client: %+v\n", err)
		t.Fail()
		return
	}
	defer cli.Close()

	for _, test := range []struct {
		name     string
		path     string
		code     int
		expected string
	}{
		{
			name: "bad id",
			path: "/some-ids",
			code: http.StatusBadRequest,
		},
		{
			name: "doesn't exist",
			path: "/123e4567-e89b-12d3-a456-426614174002",
			code: http.StatusNotFound,
		},
		{
			name:     "get first",
			path:     "/123e4567-e89b-12d3-a456-426614174000",
			code:     http.StatusOK,
			expected: "{\"email\":\"email-0\",\"name\":\"user-0\",\"verified\":true}",
		},
		{
			name:     "get second",
			path:     "/123e4567-e89b-12d3-a456-426614174001",
			code:     http.StatusOK,
			expected: "{\"email\":\"email-1\",\"name\":\"user-1\",\"verified\":true}",
		},
	} {
		t.Run(test.name, func(u *testing.T) {
			loadUsers(u, cli, []string{"123e4567-e89b-12d3-a456-426614174000", "123e4567-e89b-12d3-a456-426614174001"})

			r := server.GetRouter(cli.TargetConfig(), main.RouteBuilder)

			req, _ := http.NewRequest("GET", test.path, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != test.code {
				fmt.Printf("Expected %d: Recieved %d\n", test.code, w.Code)
				u.Fail()
			}

			if test.expected != "" && test.expected != w.Body.String() {
				fmt.Printf("Expceted %+v: Recieved %+v", test.expected, w.Body.String())
				u.Fail()
			}
		})
	}
}
