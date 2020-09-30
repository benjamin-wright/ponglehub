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

	exists, err := keycloak.HasRealm(cfg.realm)
	if err != nil {
		logrus.Fatalf("Failed checking for realm: %+v", err)
	}

	if exists {
		err = keycloak.RemoveRealm(cfg.realm)
		if err != nil {
			logrus.Fatalf("Failed to remove realm: %+v", err)
		}
	} else {
		logrus.Infof("%s realm not currently installed", cfg.realm)
	}

	err = keycloak.AddRealm(cfg.realm)
	if err != nil {
		logrus.Fatalf("Failed to add realm: %+v", err)
	}

	logrus.Info("Finished")
}
