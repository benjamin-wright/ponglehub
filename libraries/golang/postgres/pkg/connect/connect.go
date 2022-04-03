package connect

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
)

func getConnection(config *pgx.ConnConfig) *pgx.Conn {
	finished := make(chan *pgx.Conn, 1)

	go func(finished chan<- *pgx.Conn) {
		attempts := 0
		limit := 10
		var connection *pgx.Conn
		var err error
		for attempts < limit {
			attempts += 1
			connection, err = pgx.ConnectConfig(context.Background(), config)
			if err != nil {
				logrus.Debugf("Error connecting to the database: %+v", err)
			} else {
				break
			}
		}

		if connection == nil {
			logrus.Warnf("Failed to connect: %+v", err)
		}

		finished <- connection
	}(finished)

	return <-finished
}

func Connect(config ConnectConfig) (*pgx.Conn, error) {
	dbSuffix := ""
	if config.Database != "" {
		dbSuffix = "/" + config.Database
	}

	connectionString := fmt.Sprintf("postgresql://%s@%s:%d%s", config.Username, config.Host, config.Port, dbSuffix)

	logrus.Infof("Connecting to postgres with connection string: %s", connectionString)

	pgxConfig, err := pgx.ParseConfig(connectionString)
	if err != nil {
		return nil, err
	}

	conn := getConnection(pgxConfig)
	if conn == nil {
		return nil, errors.New("failed to create connection, exiting")
	}

	return conn, nil
}
