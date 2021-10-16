package migrate

import (
	"fmt"

	"ponglehub.co.uk/auth/db-init/internal/database"
	"ponglehub.co.uk/auth/db-init/internal/migrations"
	"ponglehub.co.uk/auth/db-init/pkg/types"
)

func Clean(config *types.MigrationConfig) error {
	db, err := database.Admin(config.Host, config.Port, config.CertsDir)
	if err != nil {
		return fmt.Errorf("failed to connect to db: %+v", err)
	}
	defer db.Stop()

	if err := db.DropDatabase(config.Database); err != nil {
		return fmt.Errorf("error dropping database: %+v", err)
	}

	if err := db.DropUser(config.Username); err != nil {
		return fmt.Errorf("error dropping username: %+v", err)
	}

	return nil
}

func Migrate(config *types.MigrationConfig) error {
	err := migrations.Initialize(config)
	if err != nil {
		return err
	}

	return migrations.Migrate(config)
}
