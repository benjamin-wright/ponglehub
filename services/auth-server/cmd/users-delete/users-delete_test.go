package main_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	main "ponglehub.co.uk/auth/auth-server/cmd/users-delete"
	"ponglehub.co.uk/auth/auth-server/internal/server"
	"ponglehub.co.uk/auth/auth-server/internal/testutils"
)

const TEST_DB = "test_users_delete"

func TestDeleteRoute(t *testing.T) {
	cli, err := testutils.NewClient(TEST_DB)
	if err != nil {
		fmt.Printf("Failed to create test client: %+v\n", err)
		t.Fail()
		return
	}
	defer cli.Drop()

	for _, test := range []struct {
		name string
		path string
		code int
	}{
		{name: "bad id", path: "/some-ids", code: http.StatusBadRequest},
		{name: "doesn't exist", path: "/123e4567-e89b-12d3-a456-426614174000", code: http.StatusNotFound},
	} {
		t.Run(test.name, func(u *testing.T) {
			r := server.GetRouter(TEST_DB, main.RouteBuilder)

			req, _ := http.NewRequest("DELETE", test.path, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != test.code {
				fmt.Printf("Expected %d: Recieved %d\n", test.code, w.Code)
				t.Fail()
			}
		})
	}
}
