package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"ponglehub.co.uk/auth/auth-server/internal/server"
)

func TestLoginRoute(t *testing.T) {
	r := server.GetRouter(routeBuilder)

	payload := struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Username: "bobby",
		Email:    "bob@by.com",
		Password: "pwd",
	}

	data, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(data))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	fmt.Printf("Returned status: %d\n", w.Code)

	if w.Code != http.StatusUnauthorized {
		t.Fail()
	}
}
