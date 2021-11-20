package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"

	"ponglehub.co.uk/lib/postgres/pkg/connect"
)

type AdminConn struct {
	conn *pgx.Conn
}

func NewAdminConn(cfg connect.ConnectConfig) (*AdminConn, error) {
	conn, err := connect.Connect(cfg)
	if err != nil {
		return nil, err
	}

	return &AdminConn{conn}, nil
}

func (d *AdminConn) Stop() {
	d.conn.Close(context.Background())
}

func (d *AdminConn) CreateUser(username string, password string) error {
	rows, err := d.conn.Query(context.Background(), "SHOW USERS")
	if err != nil {
		return fmt.Errorf("failed to fetch existing user: %+v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var existing string
		if err := rows.Scan(&existing, nil, nil); err != nil {
			return fmt.Errorf("failed to decode existing database user: %+v", err)
		}

		if existing == username {
			logrus.Infof("User %s already exists!", username)
			return nil
		}
	}

	logrus.Infof("Creating user %s!", username)
	if _, err := d.conn.Exec(context.Background(), "CREATE USER $1 WITH PASSWORD $2", username, password); err != nil {
		return fmt.Errorf("failed to create database user: %+v", err)
	}

	return nil
}

func (d *AdminConn) DropUser(username string) error {
	rows, err := d.conn.Query(context.Background(), "SHOW USERS")
	if err != nil {
		return fmt.Errorf("failed to fetch existing database user: %+v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var existing string
		if err := rows.Scan(&existing, nil, nil); err != nil {
			return fmt.Errorf("failed to decode existing database user: %+v", err)
		}

		if existing == username {
			rows.Close()

			logrus.Infof("Deleting user %s!", username)
			if _, err := d.conn.Exec(context.Background(), "DROP USER $1", username); err != nil {
				return fmt.Errorf("failed to drop database user: %+v", err)
			}

			return nil
		}
	}

	logrus.Infof("User %s doesn't exist!", username)
	return nil
}

func (d *AdminConn) CreateDatabase(database string) error {
	rows, err := d.conn.Query(context.Background(), "SELECT datname FROM pg_database")
	if err != nil {
		return fmt.Errorf("failed to fetch existing database: %+v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var existing string
		if err := rows.Scan(&existing); err != nil {
			return fmt.Errorf("failed to decode existing database: %+v", err)
		}

		if existing == database {
			logrus.Infof("Database %s already exists!", database)
			return nil
		}
	}

	logrus.Infof("Creating database %s!", database)
	if _, err := d.conn.Exec(context.Background(), fmt.Sprintf("CREATE DATABASE %s", database)); err != nil {
		return fmt.Errorf("failed to create database: %+v", err)
	}

	return nil
}

func (d *AdminConn) DropDatabase(database string) error {
	rows, err := d.conn.Query(context.Background(), "SELECT datname FROM pg_database")
	if err != nil {
		return fmt.Errorf("failed to fetch existing database: %+v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var existing string
		if err := rows.Scan(&existing); err != nil {
			return fmt.Errorf("failed to decode existing database: %+v", err)
		}

		if existing == database {
			rows.Close()

			logrus.Infof("Dropping database %s!", database)
			if _, err := d.conn.Exec(context.Background(), fmt.Sprintf("DROP DATABASE %s", database)); err != nil {
				return fmt.Errorf("failed to drop database: %+v", err)
			}

			return nil
		}
	}

	logrus.Infof("Database %s didn't exist!", database)
	return nil
}

func (d *AdminConn) GrantPermissions(username string, database string) error {
	query := fmt.Sprintf("GRANT ALL ON DATABASE %s TO %s", database, username)
	if _, err := d.conn.Exec(context.Background(), query); err != nil {
		return fmt.Errorf("failed to grant permissions: %+v", err)
	}

	logrus.Infof("Granted '%s' permission to read/write to '%s'!", username, database)

	return nil
}

func (d *AdminConn) RevokePermissions(username string, database string) error {
	rows, err := d.conn.Query(context.Background(), "SHOW USERS")
	if err != nil {
		return fmt.Errorf("failed to fetch existing users: %+v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var existing string
		if err := rows.Scan(&existing, nil, nil); err != nil {
			return fmt.Errorf("failed to decode existing user: %+v", err)
		}

		if existing == username {
			rows.Close()

			query := fmt.Sprintf("REVOKE ALL ON DATABASE %s FROM %s", database, username)
			if _, err := d.conn.Exec(context.Background(), query); err != nil {
				return fmt.Errorf("failed to revoke permissions: %+v", err)
			}

			logrus.Infof("Revoked '%s' permission to read/write from '%s'!", username, database)
			return nil
		}
	}

	logrus.Infof("User '%s' doesn't exist!", username)
	return nil
}
