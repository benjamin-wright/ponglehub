package migrate

import (
	"ponglehub.co.uk/lib/postgres/internal/migrations"
	"ponglehub.co.uk/lib/postgres/pkg/connect"
	"ponglehub.co.uk/lib/postgres/pkg/types"
)

// func Clean(config *types.MigrationConfig) error {
// 	db, err := database.NewAdminConn(config.AdminConfig)
// 	if err != nil {
// 		return fmt.Errorf("failed to connect to db: %+v", err)
// 	}
// 	defer db.Stop()

// 	if err := db.DropDatabase(config.TargetConfig.Database); err != nil {
// 		return fmt.Errorf("error dropping database: %+v", err)
// 	}

// 	if err := db.DropUser(config.TargetConfig.Username); err != nil {
// 		return fmt.Errorf("error dropping username: %+v", err)
// 	}

// 	return nil
// }

func Initialize(config connect.ConnectConfig, database string, username string) error {
	return migrations.Initialize(config, database, username)
}

func UnInitialize(config connect.ConnectConfig, database string, username string) error {
	return migrations.UnInitialize(config, database, username)
}

func Migrate(config connect.ConnectConfig, queries []types.Migration) error {
	return migrations.Migrate(config, queries)
}
