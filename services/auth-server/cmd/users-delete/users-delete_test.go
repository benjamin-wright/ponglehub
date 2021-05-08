package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"ponglehub.co.uk/auth/auth-server/internal/server"
)

func TestDeleteRoute(t *testing.T) {
	for _, test := range []struct {
		name string
		path string
		code int
	}{
		{name: "bad id", path: "/some-id", code: http.StatusBadRequest},
		{name: "doesn't exist", path: "/123e4567-e89b-12d3-a456-426614174000", code: http.StatusNotFound},
	} {
		t.Run(test.name, func(u *testing.T) {
			r := server.GetRouter(routeBuilder)

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
