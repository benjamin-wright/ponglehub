package main_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"

	main "ponglehub.co.uk/auth/auth-server/cmd/users-delete"
	"ponglehub.co.uk/auth/auth-server/internal/server"
	"ponglehub.co.uk/auth/auth-server/internal/testutils"
)

const TEST_DB = "test_users_delete"

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

func checkExpected(t *testing.T, cli *testutils.TestClient, expected []string) {
	ids, err := cli.ListUserIds()

	if err != nil {
		t.Fatalf("Failed to fetch user ids: %+v", err)
	}

	sort.Strings(ids)
	sort.Strings(expected)

	if len(ids) != len(expected) {
		t.Fatalf("Expected %v to equal %v", ids, expected)
	}

	for idx, _ := range ids {
		if ids[idx] != expected[idx] {
			t.Fatalf("Expected %v to equal %v", ids, expected)
		}
	}
}

func TestDeleteRoute(t *testing.T) {
	cli, err := testutils.NewClient(TEST_DB)
	if err != nil {
		fmt.Printf("Failed to create test client: %+v\n", err)
		t.Fail()
		return
	}
	defer cli.Drop()

	for _, test := range []struct {
		name     string
		path     string
		code     int
		ids      []string
		expected []string
	}{
		{name: "bad id", path: "/some-ids", code: http.StatusBadRequest},
		{
			name:     "doesn't exist",
			path:     "/123e4567-e89b-12d3-a456-426614174000",
			code:     http.StatusNotFound,
			ids:      []string{"123e4567-e89b-12d3-a456-426614174999"},
			expected: []string{"123e4567-e89b-12d3-a456-426614174999"},
		},
		{
			name: "delete successfully",
			path: "/123e4567-e89b-12d3-a456-426614174000",
			code: http.StatusNoContent,
			ids:  []string{"123e4567-e89b-12d3-a456-426614174000"},
		},
		{
			name:     "delete first",
			path:     "/123e4567-e89b-12d3-a456-426614174000",
			code:     http.StatusNoContent,
			ids:      []string{"123e4567-e89b-12d3-a456-426614174000", "123e4567-e89b-12d3-a456-426614174999"},
			expected: []string{"123e4567-e89b-12d3-a456-426614174999"},
		},
		{
			name:     "delete last",
			path:     "/123e4567-e89b-12d3-a456-426614174999",
			code:     http.StatusNoContent,
			ids:      []string{"123e4567-e89b-12d3-a456-426614174000", "123e4567-e89b-12d3-a456-426614174999"},
			expected: []string{"123e4567-e89b-12d3-a456-426614174000"},
		},
	} {
		t.Run(test.name, func(u *testing.T) {
			loadUsers(u, cli, test.ids)

			r := server.GetRouter(TEST_DB, main.RouteBuilder)

			req, _ := http.NewRequest("DELETE", test.path, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != test.code {
				fmt.Printf("Expected %d: Recieved %d\n", test.code, w.Code)
				u.Fail()
			}

			checkExpected(u, cli, test.expected)
		})
	}
}
