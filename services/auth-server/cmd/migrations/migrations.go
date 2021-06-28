package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/auth/db-init/pkg/migrate"
)

func main() {
	host, ok := os.LookupEnv("DB_HOST")
	if !ok {
		logrus.Fatalf("Failed to lookup DB_HOST env var")
	}

	user, ok := os.LookupEnv("DB_USER")
	if !ok {
		logrus.Fatalf("Failed to lookup DB_USER env var")
	}

	database, ok := os.LookupEnv("DB_NAME")
	if !ok {
		logrus.Fatal("Failed to lookup DB_NAME env var")
	}

	migrate.Migrate(&migrate.MigrationConfig{
		Host:       host,
		Port:       26257,
		Username:   user,
		Database:   database,
		Migrations: "/migrations",
	})
}
