package main_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"ponglehub.co.uk/auth-server-mock/internal/routes"
	"ponglehub.co.uk/auth-server-mock/internal/state"
)

func TestDeleteRoute(t *testing.T) {
	for _, test := range []struct {
		name string
	}{
		{name: "works"},
	} {
		t.Run(test.name, func(u *testing.T) {
			r := gin.Default()
			s := state.New()

			routes.RouteBuilder(r, s)

			req, _ := http.NewRequest("GET", "/users", nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != 200 {
				fmt.Printf("Expected %d: Recieved %d\n", 200, w.Code)
				u.Fail()
			}
		})
	}
}
