package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
)

type Database struct {
	conn  *pgx.Conn
	admin bool
}

func New(host string, port int, username string, database string) (*Database, error) {
	pgxConfig, err := pgx.ParseConfig(fmt.Sprintf("postgres://%s@%s:%d/%s", username, host, port, database))
	if err != nil {
		return nil, err
	}

	conn := getConnection(pgxConfig)
	if conn == nil {
		return nil, errors.New("failed to create connection, exiting")
	}

	return &Database{
		conn:  conn,
		admin: false,
	}, nil
}

func Admin(host string, port int) (*Database, error) {
	pgxConfig, err := pgx.ParseConfig(fmt.Sprintf("postgres://root@%s:%d", host, port))
	if err != nil {
		return nil, err
	}

	conn := getConnection(pgxConfig)
	if conn == nil {
		return nil, errors.New("failed to create connection, exiting")
	}

	return &Database{
		conn:  conn,
		admin: true,
	}, nil
}

func getConnection(config *pgx.ConnConfig) *pgx.Conn {
	finished := make(chan *pgx.Conn, 1)

	go func(finished chan<- *pgx.Conn) {
		attempts := 0
		limit := 10
		var connection *pgx.Conn
		var err error
		for attempts < limit {
			connection, err = pgx.ConnectConfig(context.Background(), config)
			if err != nil {
				logrus.Warnf("error connecting to the database: %+v", err)
			} else {
				break
			}
		}

		finished <- connection
	}(finished)

	return <-finished
}

func (d *Database) Stop() {
	d.conn.Close(context.Background())
}

func (d *Database) CreateUser(username string) error {
	if !d.admin {
		return errors.New("cannot call CreateUser on non-admin connection")
	}

	rows, err := d.conn.Query(context.Background(), "SHOW USERS")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var existing string
		if err := rows.Scan(&existing, nil, nil); err != nil {
			return err
		}

		if existing == username {
			logrus.Infof("User %s already exists!", username)
			return nil
		}
	}

	logrus.Infof("Creating user %s!", username)
	if _, err := d.conn.Exec(context.Background(), "CREATE USER $1", username); err != nil {
		return err
	}

	return nil
}

func (d *Database) CreateDatabase(database string) error {
	if !d.admin {
		return errors.New("cannot call CreateDatabase on non-admin connection")
	}

	rows, err := d.conn.Query(context.Background(), "SELECT datname FROM pg_database")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var existing string
		if err := rows.Scan(&existing); err != nil {
			return err
		}

		if existing == database {
			logrus.Infof("Database %s already exists!", database)
			return nil
		}
	}

	logrus.Infof("Creating database %s!", database)
	if _, err := d.conn.Exec(context.Background(), fmt.Sprintf("CREATE DATABASE %s", database)); err != nil {
		return err
	}

	return nil
}

func (d *Database) GrantPermissions(username string, database string) error {
	if !d.admin {
		return errors.New("cannot call GrantPermissions on non-admin connection")
	}

	query := fmt.Sprintf("GRANT ALL ON DATABASE %s TO %s", database, username)
	if _, err := d.conn.Exec(context.Background(), query); err != nil {
		return err
	}

	logrus.Infof("Granted '%s' permission to read/write to '%s'!", username, database)

	return nil
}

func (d *Database) EnsureMigrationTable() error {
	_, err := d.conn.Exec(
		context.TODO(),
		`
			BEGIN;

			SAVEPOINT migration_restart;

			CREATE TABLE IF NOT EXISTS migrations (
				id INT PRIMARY KEY NOT NULL UNIQUE
			);

			RELEASE SAVEPOINT migration_restart;

			COMMIT;
		`,
	)

	return err
}
