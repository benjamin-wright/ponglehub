package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/auth/auth-server/internal/migrations"
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

	password, ok := os.LookupEnv("DB_PASS")
	if !ok {
		logrus.Fatal("Enrivonment Variable DB_PASS not found")
	}

	certsDir, ok := os.LookupEnv("DB_CERTS")
	if !ok {
		logrus.Fatal("Enrivonment Variable DB_CERTS not found")
	}

	if err := migrations.Migrate(host, user, password, database, certsDir); err != nil {
		logrus.Fatalf("Failed to run migrations: %+v", err)
	}
}
