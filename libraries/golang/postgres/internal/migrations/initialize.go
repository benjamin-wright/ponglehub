package migrations

import (
	"fmt"

	"ponglehub.co.uk/lib/postgres/internal/database"
	"ponglehub.co.uk/lib/postgres/pkg/connect"
)

func Initialize(config connect.ConnectConfig, dbName string, username string) error {
	db, err := database.NewAdminConn(config)
	if err != nil {
		return fmt.Errorf("failed to connect to db: %+v", err)
	}
	defer db.Stop()

	if err := db.CreateUser(username); err != nil {
		return fmt.Errorf("error creating user: %+v", err)
	}

	if err := db.CreateDatabase(dbName); err != nil {
		return fmt.Errorf("error creating user: %+v", err)
	}

	if err := db.GrantPermissions(username, dbName); err != nil {
		return fmt.Errorf("error granting permissions: %+v", err)
	}

	return nil
}

func UnInitialize(config connect.ConnectConfig, dbName string, username string) error {
	db, err := database.NewAdminConn(config)
	if err != nil {
		return fmt.Errorf("failed to connect to db: %+v", err)
	}
	defer db.Stop()

	if err := db.RevokePermissions(username, dbName); err != nil {
		return fmt.Errorf("error revoking permissions: %+v", err)
	}

	if err := db.DropDatabase(dbName); err != nil {
		return fmt.Errorf("error dropping user: %+v", err)
	}

	if err := db.DropUser(username); err != nil {
		return fmt.Errorf("error dropping user: %+v", err)
	}

	return nil
}
