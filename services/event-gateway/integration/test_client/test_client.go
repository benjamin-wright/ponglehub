package test_client

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/cookiejar"
	"testing"

	"github.com/stretchr/testify/assert"
)

func noErr(t *testing.T, err error) {
	if err != nil {
		assert.NoError(t, err)
		t.FailNow()
	}
}

type TestClient struct {
	client *http.Client
	jar    http.CookieJar
}

func New(t *testing.T) *TestClient {
	jar, err := cookiejar.New(nil)
	noErr(t, err)

	return &TestClient{
		client: &http.Client{
			Jar: jar,
		},
		jar: jar,
	}
}

func (t *TestClient) CookieJar() http.CookieJar {
	return t.jar
}

func (t *TestClient) Post(u *testing.T, url string, data map[string]string) *http.Response {
	json_data, err := json.Marshal(data)
	noErr(u, err)

	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(json_data),
	)
	noErr(u, err)

	req.Header.Add("Content-Type", "application/json")

	res, err := t.client.Do(req)
	noErr(u, err)

	return res
}
