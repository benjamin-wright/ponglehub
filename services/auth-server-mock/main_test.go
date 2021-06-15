package main_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"ponglehub.co.uk/auth-server-mock/internal/routes"
	"ponglehub.co.uk/auth-server-mock/internal/state"
)

func toAnon(users []state.User) []map[string]interface{} {
	output := []map[string]interface{}{}

	for _, user := range users {
		output = append(output, map[string]interface{}{
			"id":       user.ID,
			"name":     user.Name,
			"email":    user.Email,
			"verified": user.Verified,
		})
	}

	return output
}

func TestListRoute(t *testing.T) {
	user1 := state.User{ID: "some-id", Name: "ben", Email: "some@email.com", Password: "hashed", Verified: false}
	user2 := state.User{ID: "some-other-id", Name: "geoss", Email: "other@email.com", Password: "alsohashed", Verified: false}

	for _, test := range []struct {
		name     string
		existing []state.User
	}{
		{name: "empty result", existing: []state.User{}},
		{name: "one element", existing: []state.User{user1}},
		{name: "two elements", existing: []state.User{user1, user2}},
	} {
		t.Run(test.name, func(u *testing.T) {
			r := gin.Default()
			s := state.New()
			s.Users = test.existing

			routes.RouteBuilder(r, s)

			req, _ := http.NewRequest("GET", "/users", nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != 200 {
				fmt.Printf("Expected %d: Recieved %d\n", 200, w.Code)
				u.FailNow()
			}

			var response []map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				fmt.Printf("Failed to parse json: %s", err.Error())
				u.FailNow()
			}

			assert.EqualValues(u, toAnon(test.existing), response, "Expected list results to be the same")
		})
	}
}
