package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/sirupsen/logrus"
)

// KeycloakAPI api for interacting with Keycloak
type KeycloakAPI struct {
	url          string
	accessToken  string
	refreshToken string
}

// New create a new API instance
func New(authURL string, username string, password string) (*KeycloakAPI, error) {
	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)
	data.Set("grant_type", "password")
	data.Set("client_id", "admin-cli")

	resp, err := http.Post(
		fmt.Sprintf("%s/auth/realms/master/protocol/openid-connect/token", authURL),
		"application/x-www-form-urlencoded",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("Recieved non-200 status code: %d - %s", resp.StatusCode, string(body))
	}

	type authResponse struct {
		AccessToken      string `json:"access_token"`
		ExpiresIn        int    `json:"expires_in"`
		RefreshExpiresIn int    `json:"refresh_expires_in"`
		RefreshToken     string `json:"refresh_token"`
		TokenType        string `json:"token_type"`
		NotBeforePolicy  int    `json:"not_before_policy"`
		SessionState     string `json:"session_state"`
		Scope            string `json:"scope"`
	}
	var parsed authResponse
	err = json.Unmarshal(body, &parsed)
	if err != nil {
		return nil, err
	}

	return &KeycloakAPI{
		url:          authURL,
		accessToken:  parsed.AccessToken,
		refreshToken: parsed.RefreshToken,
	}, nil
}

func (k *KeycloakAPI) get(path string, response interface{}) (int, error) {
	url := fmt.Sprintf("%s/auth/admin/realms%s", k.url, path)
	logrus.Debugf("Submitting get request to '%s'", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+k.accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return resp.StatusCode, fmt.Errorf("[%d] %s", resp.StatusCode, string(resBody))
	}

	if response != nil {
		err = json.Unmarshal(resBody, &response)
		if err != nil {
			return resp.StatusCode, err
		}
	}

	return resp.StatusCode, nil
}

func (k *KeycloakAPI) post(path string, body []byte, response interface{}) (int, error) {
	url := fmt.Sprintf("%s/auth/admin/realms%s", k.url, path)
	logrus.Debugf("Submitting post request to '%s'", url)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+k.accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return resp.StatusCode, fmt.Errorf("[%d] %s", resp.StatusCode, string(resBody))
	}

	if response != nil {
		err = json.Unmarshal(resBody, &response)
		if err != nil {
			return resp.StatusCode, err
		}
	}

	return resp.StatusCode, nil
}

func (k *KeycloakAPI) delete(path string, response interface{}) (int, error) {
	url := fmt.Sprintf("%s/auth/admin/realms%s", k.url, path)
	logrus.Debugf("Submitting delete request to '%s'", url)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+k.accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return resp.StatusCode, fmt.Errorf("[%d] %s", resp.StatusCode, string(resBody))
	}

	if response != nil {
		err = json.Unmarshal(resBody, &response)
		if err != nil {
			return resp.StatusCode, err
		}
	}

	return resp.StatusCode, nil
}
