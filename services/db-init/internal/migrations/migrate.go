package migrations

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/auth/db-init/internal/database"
	"ponglehub.co.uk/auth/db-init/pkg/types"
)

func Migrate(config *types.MigrationConfig) error {
	db, err := database.New(config.Host, config.Port, config.Username, config.Password, config.Database, config.CertsDir)
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
