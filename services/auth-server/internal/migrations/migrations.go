package migrations

import (
	"ponglehub.co.uk/lib/postgres/pkg/connect"
	m "ponglehub.co.uk/lib/postgres/pkg/migrate"
	"ponglehub.co.uk/lib/postgres/pkg/types"
)

func Migrate(targetConfig connect.ConnectConfig, adminConfig connect.ConnectConfig) error {
	return m.Migrate(&types.MigrationConfig{
		TargetConfig: targetConfig,
		AdminConfig:  adminConfig,
		Migrations: []types.Migration{
			{Query: `
				CREATE TABLE users (
					id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
					name VARCHAR(100) NOT NULL UNIQUE,
					email VARCHAR(100) NOT NULL UNIQUE,
					password VARCHAR(100),
					verified BOOLEAN NOT NULL
				);
			`},
		},
	})
}
