package migrate_test

import (
	"io"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"ponglehub.co.uk/auth/db-init/pkg/migrate"
)

func TestMigrations(t *testing.T) {
	logrus.SetOutput(io.Discard)

	t.Run("TestCase", func(u *testing.T) {
		config := migrate.MigrationConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     26257,
			Username: os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			CertsDir: os.Getenv("DB_CERTS"),
			Database: "test_db",
			Migrations: []migrate.Migration{
				{
					Query: `
						BEGIN;

						SAVEPOINT migration_1_restart;
						
						DROP TABLE IF EXISTS test_users;
						
						CREATE TABLE test_users (
							id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
							name VARCHAR(100) NOT NULL UNIQUE,
							email VARCHAR(100) NOT NULL UNIQUE,
							password VARCHAR(100),
							verified BOOLEAN NOT NULL
						);
						
						RELEASE SAVEPOINT migration_1_restart;
						
						COMMIT;
					`,
				},
			},
		}

		err := migrate.Clean(&config)
		assert.NoError(u, err)

		err = migrate.Migrate(&config)
		assert.NoError(u, err)
	})
}
