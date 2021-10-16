package migrations

import (
	"fmt"

	"ponglehub.co.uk/auth/db-init/internal/database"
	"ponglehub.co.uk/auth/db-init/pkg/types"
)

func Initialize(config *types.MigrationConfig) error {
	db, err := database.Admin(config.Host, config.Port, config.CertsDir)
	if err != nil {
		return fmt.Errorf("failed to connect to db: %+v", err)
	}
	defer db.Stop()

	if err := db.CreateUser(config.Username, config.Password); err != nil {
		return fmt.Errorf("error creating user: %+v", err)
	}

	if err := db.CreateDatabase(config.Database); err != nil {
		return fmt.Errorf("error creating user: %+v", err)
	}

	if err := db.GrantPermissions(config.Username, config.Database); err != nil {
		return fmt.Errorf("error granting permissions: %+v", err)
	}

	return nil
}
