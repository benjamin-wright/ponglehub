package main

import (
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/keycloak-init/internal/api"
)

func main() {
	logrus.Info("Starting...")

	cfg, err := newConfig()

	if err != nil {
		logrus.Fatalf("Failed to load config: %+v", err)
	}

	cfg.print()

	keycloak, err := api.New(cfg.url, cfg.username, cfg.password)
	if err != nil {
		logrus.Fatalf("Failed to create api instance: %+v", err)
	}

	realm := api.Realm{
		Name:    cfg.realm,
		Display: cfg.realm,
		SMTPServer: &api.SMTPServer{
			User:               "admin@ponglehub.co.uk",
			Password:           "******",
			ReplyToDisplayName: "Ponglehub",
			StartTLS:           false,
			Auth:               true,
			Port:               465,
			Host:               "mail.gandi.net",
			ReplyTo:            "admin@ponglehub.co.uk",
			From:               "admin@ponglehub.co.uk",
			FromDisplayName:    "Ponglehub",
			SSL:                true,
		},
	}

	exists, err := keycloak.HasRealm(realm)
	if err != nil {
		logrus.Fatalf("Failed checking for realm: %+v", err)
	}

	if exists {
		err = keycloak.RemoveRealm(realm)
		if err != nil {
			logrus.Fatalf("Failed to remove realm: %+v", err)
		}
	} else {
		logrus.Infof("%s realm not currently installed", cfg.realm)
	}

	err = keycloak.AddRealm(realm)
	if err != nil {
		logrus.Fatalf("Failed to add realm: %+v", err)
	}

	logrus.Info("Finished")
}
