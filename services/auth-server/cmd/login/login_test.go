package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"ponglehub.co.uk/auth/auth-server/internal/server"
	"ponglehub.co.uk/auth/auth-server/internal/testutils"
)

const TEST_DB = "test_login"

func TestLoginRoute(t *testing.T) {
	cli, err := testutils.NewClient(TEST_DB)
	if err != nil {
		fmt.Printf("Failed to create test client: %+v\n", err)
		t.Fail()
		return
	}
	defer cli.Close()

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
			if err := cli.Reset(); err != nil {
				u.Fatalf("Error clearing database: %+v", err)
			}

			err = cli.AddUser("123e4567-e89b-12d3-a456-426614174002", "johnny", "johnny@place.com", "some-pass", true)
			if err != nil {
				fmt.Printf("Failed to add user: %+v\n", err)
				u.Fail()
				return
			}

			r := server.GetRouter(cli.TargetConfig(), routeBuilder)

			payload := struct {
				Email    string `json:"email"`
				Password string `json:"password"`
			}{
				Email:    test.email,
				Password: test.password,
			}

			data, _ := json.Marshal(payload)

			req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(data))
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != test.code {
				fmt.Printf("Expected %d: Recieved %d\n", test.code, w.Code)
				u.Fail()
			}
		})
	}
}
