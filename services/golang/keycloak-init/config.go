package main

import (
	"fmt"

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

func (c *config) String() string {
	return fmt.Sprintf(
		"{url: %s, realm: %s, username: %s, email: %s, host: %s, port: %d, from: %s}",
		c.URL,
		c.Realm,
		c.Username,
		c.SMTPEmail,
		c.SMTPHost,
		c.SMTPPort,
		c.SMTPFrom,
	)
}
