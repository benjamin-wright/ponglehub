package migrate_test

import (
	"io"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"ponglehub.co.uk/auth/db-init/internal/database"
	"ponglehub.co.uk/auth/db-init/pkg/migrate"
)

func TestMigrations(t *testing.T) {
	logrus.SetOutput(io.Discard)

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
	assert.NoError(t, err)

	err = migrate.Migrate(&config)
	assert.NoError(t, err)

	t.Run("has tables", func(u *testing.T) {
		db, err := database.New(os.Getenv("DB_HOST"), 26257, os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), "test_db", os.Getenv("DB_CERTS"))
		if err != nil {
			assert.NoError(u, err)
			assert.FailNow(u, "error connecting to db")
		}

		tables, err := db.GetTables()
		if err != nil {
			assert.NoError(u, err)
			assert.FailNow(u, "error fetching tables")
		}

		assert.Contains(u, tables, "migrations")
		assert.Contains(u, tables, "test_users")
	})
}
