package testutils

import (
	"context"
	"errors"
	"fmt"
	"os"

	pgx "github.com/jackc/pgx/v4"
)

type TestClient struct {
	conn     *pgx.Conn
	database string
}

// New - Create a new AuthClient instance
func NewClient(database string) (*TestClient, error) {
	if err := migrate(database); err != nil {
		return nil, fmt.Errorf("Failed to migrate database: %+v", err)
	}

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

	return &TestClient{conn: conn, database: database}, nil
}

func (c *TestClient) Drop() error {
	c.conn.Close(context.Background())

	host, ok := os.LookupEnv("DB_HOST")
	if !ok {
		return errors.New("DB_HOST env var must be provided")
	}

	config, err := pgx.ParseConfig(fmt.Sprintf("postgres://root@%s:26257", host))
	if err != nil {
		return err
	}

	conn, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		return fmt.Errorf("error connecting to the database: %+v", err)
	}
	defer conn.Close(context.Background())

	_, err = conn.Exec(
		context.TODO(),
		fmt.Sprintf("DROP DATABASE %s;", c.database),
	)

	return err
}

func migrate(database string) error {
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
				fmt.Printf("error connecting to the database: %+v\n", err)
			} else {
				break
			}
		}

		finished <- connection
		return
	}(finished)

	conn := <-finished
	if conn == nil {
		return errors.New("Failed to create connection, exiting.")
	}
	defer conn.Close(context.Background())

	_, err = conn.Exec(
		context.TODO(),
		fmt.Sprintf(`
			CREATE DATABASE %s;

			BEGIN;

			SAVEPOINT migration_1_restart;

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
