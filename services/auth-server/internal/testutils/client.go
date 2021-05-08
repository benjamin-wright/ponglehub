package testutils

import (
	"context"
	"errors"
	"fmt"
	"os"

	pgx "github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
)

type TestClient struct {
	conn  *pgx.Conn
	table string
}

// New - Create a new AuthClient instance
func NewClient(database string) (*TestClient, error) {
	host, ok := os.LookupEnv("DB_HOST")
	if !ok {
		return nil, errors.New("DB_HOST env var must be provided")
	}

	pgxConfig, err := pgx.ParseConfig(fmt.Sprintf("postgres://root@%s:26257/%s", host, database))
	if err != nil {
		return nil, err
	}

	conn, err := pgx.ConnectConfig(context.TODO(), pgxConfig)
	if err != nil {
		return nil, err
	}

	return &TestClient{conn: conn}, nil
}

func Migrate(database string) error {
	host, ok := os.LookupEnv("DB_HOST")
	if !ok {
		return errors.New("DB_HOST env var must be provided")
	}

	config, err := pgx.ParseConfig(fmt.Sprintf("postgres://root@%s:26257", host))
	if err != nil {
		return err
	}

	finished := make(chan *pgx.Conn, 1)

	go func(finished chan<- *pgx.Conn) {
		attempts := 0
		limit := 10
		var connection *pgx.Conn
		for attempts < limit {
			connection, err = pgx.ConnectConfig(context.Background(), config)
			if err != nil {
				logrus.Warnf("error connecting to the database: %+v", err)
			} else {
				break
			}
		}

		finished <- connection
		return
	}(finished)

	conn := <-finished
	if conn == nil {
		logrus.Fatalf("Failed to create connection, exiting.")
	}
	defer conn.Close(context.Background())

	conn.Exec(
		context.TODO(),
		fmt.Sprintf(`
			CREATE DATABASE %s;
		`, database),
	)

	fmt.Println("Created database (or not)")

	_, err = conn.Exec(
		context.TODO(),
		fmt.Sprintf(`
			BEGIN;

			SAVEPOINT migration_1_restart;

			DROP TABLE IF EXISTS %s.users;

			CREATE TABLE %s.users (
					id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
					name VARCHAR(100) NOT NULL UNIQUE,
					email VARCHAR(100) NOT NULL UNIQUE,
					password VARCHAR(100),
					verified BOOLEAN NOT NULL
			);

			RELEASE SAVEPOINT migration_1_restart;

			COMMIT;
		`, database, database),
	)

	return err
}
