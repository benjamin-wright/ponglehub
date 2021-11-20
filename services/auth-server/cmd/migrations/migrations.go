package main

import (
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/auth/auth-server/internal/migrations"
	"ponglehub.co.uk/lib/postgres/pkg/connect"
)

func main() {
	targetCfg, err := connect.ConfigFromEnv()
	if err != nil {
		logrus.Fatalf("Failed to load target config from environment: %+v", err)
	}

	adminCfg, err := connect.AdminFromEnv()
	if err != nil {
		logrus.Fatalf("Failed to load admin config from environment: %+v", err)
	}

	if err := migrations.Migrate(targetCfg, adminCfg); err != nil {
		logrus.Fatalf("Failed to run migrations: %+v", err)
	}
}
