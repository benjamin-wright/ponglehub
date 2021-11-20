package migrate

import (
	"fmt"

	"ponglehub.co.uk/lib/postgres/internal/database"
	"ponglehub.co.uk/lib/postgres/internal/migrations"
	"ponglehub.co.uk/lib/postgres/pkg/types"
)

func Clean(config *types.MigrationConfig) error {
	db, err := database.NewAdminConn(config.AdminConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to db: %+v", err)
	}
	defer db.Stop()

	if err := db.DropDatabase(config.TargetConfig.Database); err != nil {
		return fmt.Errorf("error dropping database: %+v", err)
	}

	if err := db.DropUser(config.TargetConfig.Username); err != nil {
		return fmt.Errorf("error dropping username: %+v", err)
	}

	return nil
}

func Migrate(config *types.MigrationConfig) error {
	err := migrations.Initialize(config)
	if err != nil {
		return err
	}

	return migrations.Migrate(config.TargetConfig, config.Migrations)
}
