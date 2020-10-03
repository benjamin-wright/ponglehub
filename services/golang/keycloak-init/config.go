package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/envreader"
)

type config struct {
	URL          string `env:"KEYCLOAK_INIT_URL"`
	Realm        string `env:"KEYCLOAK_INIT_REALM"`
	Username     string `env:"KEYCLOAK_INIT_USER"`
	Password     string `env:"KEYCLOAK_INIT_PASSWORD"`
	SMTPEmail    string `env:"KEYCLOAK_SMTP_EMAIL"`
	SMTPPassword string `env:"KEYCLOAK_SMTP_PASSWORD"`
	SMTPHost     string `env:"KEYCLOAK_SMTP_HOST"`
	SMTPPort     int    `env:"KEYCLOAK_SMTP_PORT"`
	SMTPFrom     string `env:"KEYCLOAK_SMTP_FROM"`
}

func newConfig() (*config, error) {
	cfg := &config{}

	err := envreader.Load(cfg)
	if err != nil {
		return nil, fmt.Errorf("Error loading config from environment: %+v", err)
	}

	return cfg, nil
}

func (c *config) print() {
	logrus.Infof("Config:\n - url: %s\n - realm: %s\n - username: %s\n - password: %t", c.URL, c.Realm, c.Username, c.Password != "")
}
