package testutils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"ponglehub.co.uk/auth/auth-server/internal/client"
	"ponglehub.co.uk/auth/auth-server/internal/migrations"
	"ponglehub.co.uk/lib/postgres/pkg/connect"
)

func Client(t *testing.T, database string) *client.PostgresClient {
	targetCfg, err := connect.ConfigFromEnv()
	if err != nil {
		assert.FailNow(t, "Failed to load target config from environment: %+v", err)
	}
	targetCfg.Database = database

	adminCfg, err := connect.AdminFromEnv()
	if err != nil {
		assert.FailNow(t, "Failed to load admin config from environment: %+v", err)
	}

	err = migrations.Migrate(targetCfg, adminCfg)
	if err != nil {
		assert.FailNow(t, "Failed to run migrations: %+v", err)
	}

	cli, err := client.NewPostgresClient(context.Background(), targetCfg)
	if err != nil {
		assert.FailNow(t, "Failed to connect to database: %+v", err)
	}

	return cli
}

func Hash(t *testing.T, password string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		assert.FailNow(t, "Error hashing user password: %+v", err)
	}

	return string(hashed)
}
