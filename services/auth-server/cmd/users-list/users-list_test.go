package main_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	main "ponglehub.co.uk/auth/auth-server/cmd/users-list"
	"ponglehub.co.uk/auth/auth-server/internal/server"
	"ponglehub.co.uk/auth/auth-server/internal/testutils"
)

const TEST_DB = "test_users_list"

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

func TestListRoute(t *testing.T) {
	cli, err := testutils.NewClient(TEST_DB)
	if err != nil {
		fmt.Printf("Failed to create test client: %+v\n", err)
		t.Fail()
		return
	}
	defer cli.Close()

	for _, test := range []struct {
		name     string
		ids      []string
		expected []string
	}{
		{
			name:     "empty database",
			ids:      []string{},
			expected: []string{},
		},
		{
			name:     "one user",
			ids:      []string{"123e4567-e89b-12d3-a456-426614174000"},
			expected: []string{`{"id":"123e4567-e89b-12d3-a456-426614174000","name":"user-0","email":"email-0","verified":true}`},
		},
		{
			name: "two users",
			ids:  []string{"123e4567-e89b-12d3-a456-426614174000", "123e4567-e89b-12d3-a456-426614174001"},
			expected: []string{
				`{"id":"123e4567-e89b-12d3-a456-426614174000","name":"user-0","email":"email-0","verified":true}`,
				`{"id":"123e4567-e89b-12d3-a456-426614174001","name":"user-1","email":"email-1","verified":true}`,
			},
		},
		{
			name: "three users",
			ids:  []string{"123e4567-e89b-12d3-a456-426614174000", "123e4567-e89b-12d3-a456-426614174001", "123e4567-e89b-12d3-a456-426614174002"},
			expected: []string{
				`{"id":"123e4567-e89b-12d3-a456-426614174000","name":"user-0","email":"email-0","verified":true}`,
				`{"id":"123e4567-e89b-12d3-a456-426614174001","name":"user-1","email":"email-1","verified":true}`,
				`{"id":"123e4567-e89b-12d3-a456-426614174002","name":"user-2","email":"email-2","verified":true}`,
			},
		},
	} {
		t.Run(test.name, func(u *testing.T) {
			loadUsers(u, cli, test.ids)

			r := server.GetRouter(cli.TargetConfig(), main.RouteBuilder)

			req, _ := http.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				fmt.Printf("Expected %d: Recieved %d\n", http.StatusOK, w.Code)
				u.Fail()
			}

			expected := fmt.Sprintf("[%s]", strings.Join(test.expected, ","))

			if expected != w.Body.String() {
				fmt.Printf("Expected %+v: Recieved %+v", expected, w.Body.String())
				u.Fail()
			}
		})
	}
}
