package migrate

import (
	"fmt"

	"ponglehub.co.uk/auth/db-init/internal/database"
	"ponglehub.co.uk/auth/db-init/internal/migrations"
)

type MigrationConfig struct {
	Host       string
	Port       int
	Username   string
	Database   string
	Migrations string
}

func Migrate(config *MigrationConfig) error {
	err := initialize(config)
	if err != nil {
		return err
	}

	db, err := database.New(config.Host, config.Port, config.Username, config.Database)
	if err != nil {
		return fmt.Errorf("failed to connect to db: %+v", err)
	}

	if err = db.EnsureMigrationTable(); err != nil {
		return fmt.Errorf("failed to ensure migration table: %+v", err)
	}

	files, err := migrations.Load(config.Migrations)
	if err != nil {
		return err
	}

	return fmt.Errorf("%+v", files)
}

func initialize(config *MigrationConfig) error {
	db, err := database.Admin(config.Host, config.Port)
	if err != nil {
		return fmt.Errorf("failed to connect to db: %+v", err)
	}
	defer db.Stop()

	if err := db.CreateUser(config.Username); err != nil {
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
