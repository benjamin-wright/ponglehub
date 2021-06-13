package main_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeleteRoute(t *testing.T) {
	r := gin.Default()

	req, _ := http.NewRequest("DELETE", test.path, nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != test.code {
		fmt.Printf("Expected %d: Recieved %d\n", test.code, w.Code)
		u.Fail()
	}
}
