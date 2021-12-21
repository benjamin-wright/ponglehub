package main

import (
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/auth/auth-server/internal/migrations"
	"ponglehub.co.uk/lib/postgres/pkg/connect"
)

func main() {
	config, err := connect.ConfigFromEnv()
	if err != nil {
		logrus.Fatalf("Failed to load target config from environment: %+v", err)
	}

	if err := migrations.Migrate(config); err != nil {
		logrus.Fatalf("Failed to run migrations: %+v", err)
	}
}
