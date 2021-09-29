package migrate

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/auth/db-init/internal/database"
)

type Migration struct {
	Query string
}

type MigrationConfig struct {
	Host       string
	Port       int
	Username   string
	Database   string
	Migrations []Migration
}

func Clean(config *MigrationConfig) error {
	db, err := database.Admin(config.Host, config.Port)
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

	for id, migration := range config.Migrations {
		if db.HasMigration(id) {
			logrus.Infof("Migration %d already done: skipping", id)
			continue
		}

		logrus.Infof("Migration %d running...", id)
		if err := db.RunMigration(migration.Query); err != nil {
			return err
		}

		if err := db.AddMigration(id); err != nil {
			return err
		}
	}

	return nil
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
