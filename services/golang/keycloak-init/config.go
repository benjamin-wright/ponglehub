package main

import (
	"errors"
	"os"

	"github.com/sirupsen/logrus"
)

type config struct {
	url      string
	realm    string
	username string
	password string
}

func newConfig() (*config, error) {
	cfg := config{}

	if url, ok := os.LookupEnv("KEYCLOAK_INIT_URL"); ok {
		cfg.url = url
	} else {
		return nil, errors.New("Value required for KEYCLOAK_INIT_URL")
	}

	if realm, ok := os.LookupEnv("KEYCLOAK_INIT_REALM"); ok {
		cfg.realm = realm
	} else {
		return nil, errors.New("Value required for KEYCLOAK_INIT_REALM")
	}

	if username, ok := os.LookupEnv("KEYCLOAK_INIT_USER"); ok {
		cfg.username = username
	} else {
		return nil, errors.New("Value required for KEYCLOAK_INIT_USER")
	}

	if password, ok := os.LookupEnv("KEYCLOAK_INIT_PASSWORD"); ok {
		cfg.password = password
	} else {
		return nil, errors.New("Value required for KEYCLOAK_INIT_PASSWORD")
	}

	return &cfg, nil
}

func (c *config) print() {
	logrus.Infof("Config:\n - url: %s\n - realm: %s\n - username: %s\n - password: %t", c.url, c.realm, c.username, c.password != "")
}
