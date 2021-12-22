package testutils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"ponglehub.co.uk/auth/auth-server/internal/client"
	"ponglehub.co.uk/lib/postgres/pkg/connect"
)

func Client(t *testing.T) *client.PostgresClient {
	cfg, err := connect.ConfigFromEnv()
	if err != nil {
		assert.FailNow(t, "Failed to load postgres config from environment: %+v", err)
	}

	cli, err := client.NewPostgresClient(context.Background(), cfg)
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
