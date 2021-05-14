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
	defer cli.Drop()

	err = cli.AddUser("123e4567-e89b-12d3-a456-426614174002", "johnny", "johnny@place.com", "some-pass", true)
	if err != nil {
		fmt.Printf("Failed to add user: %+v\n", err)
		t.Fail()
		return
	}

	r := server.GetRouter(TEST_DB, routeBuilder)

	payload := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Email:    "johnny@place.com",
		Password: "some-pass",
	}

	data, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(data))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		fmt.Printf("Expected %d: Recieved %d\n", http.StatusOK, w.Code)
		t.Fail()
	}
}
