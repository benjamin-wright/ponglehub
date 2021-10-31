package migrations

import (
	"fmt"

	"ponglehub.co.uk/auth/postgres/internal/database"
	"ponglehub.co.uk/auth/postgres/pkg/types"
)

func Initialize(config *types.MigrationConfig) error {
	db, err := database.NewAdminConn(config.AdminConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to db: %+v", err)
	}
	defer db.Stop()

	if err := db.CreateUser(config.TargetConfig.Username, config.TargetConfig.Password); err != nil {
		return fmt.Errorf("error creating user: %+v", err)
	}

	if err := db.CreateDatabase(config.TargetConfig.Database); err != nil {
		return fmt.Errorf("error creating user: %+v", err)
	}

	if err := db.GrantPermissions(config.TargetConfig.Username, config.TargetConfig.Database); err != nil {
		return fmt.Errorf("error granting permissions: %+v", err)
	}

	return nil
}
