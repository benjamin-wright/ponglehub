package migrations

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/auth/postgres/internal/database"
	"ponglehub.co.uk/auth/postgres/pkg/connect"
	"ponglehub.co.uk/auth/postgres/pkg/types"
)

func Migrate(config connect.ConnectConfig, migrations []types.Migration) error {
	db, err := database.NewMigrationConn(config)
	if err != nil {
		return fmt.Errorf("failed to connect to db: %+v", err)
	}

	if err = db.EnsureMigrationTable(); err != nil {
		return fmt.Errorf("failed to ensure migration table: %+v", err)
	}

	for id, migration := range migrations {
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
