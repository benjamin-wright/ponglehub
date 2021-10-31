package migrations

import (
	m "ponglehub.co.uk/auth/db-init/pkg/migrate"
	"ponglehub.co.uk/auth/db-init/pkg/types"
)

func Migrate(host string, user string, password string, database string, certsDir string) error {
	return m.Migrate(&types.MigrationConfig{
		Host:     host,
		Port:     26257,
		Username: user,
		Password: password,
		Database: database,
		CertsDir: certsDir,
		Migrations: []types.Migration{
			{Query: `
				BEGIN;

				SAVEPOINT migration_1_restart;
				
				DROP TABLE IF EXISTS users;
				
				CREATE TABLE users (
					id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
					name VARCHAR(100) NOT NULL UNIQUE,
					email VARCHAR(100) NOT NULL UNIQUE,
					password VARCHAR(100),
					verified BOOLEAN NOT NULL
				);
				
				RELEASE SAVEPOINT migration_1_restart;
				
				COMMIT;
			`},
		},
	})
}
