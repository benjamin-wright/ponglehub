package main

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type migrationSet struct {
	User     string `yaml:"user"`
	Database string `yaml:"database"`
}

type conf struct {
	Host string         `yaml:"host"`
	Port int16          `yaml:"port"`
	Data []migrationSet `yaml:"data"`
}

func main() {
	var c conf

	yamlFile, err := ioutil.ReadFile("/config.yaml")
	if err != nil {
		logrus.Fatalf("yamlFile.Get err   %+v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		logrus.Fatalf("Unmarshal: %v", err)
	}

	config, err := pgx.ParseConfig(fmt.Sprintf("postgres://root@%s:%d", c.Host, c.Port))
	if err != nil {
		logrus.Fatalf("error configuring the database: %+v", err)
	}

	// Connect to the "bank" database.
	conn, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		logrus.Fatalf("error connecting to the database: %+v", err)
	}
	defer conn.Close(context.Background())

	for _, set := range c.Data {
		if err := createUser(conn, set.User); err != nil {
			logrus.Fatalf("Error creating user: %+v", err)
		}

		if err := createDatabase(conn, set.Database); err != nil {
			logrus.Fatalf("Error creating user: %+v", err)
		}

		if err := grantPermissions(conn, set.User, set.Database); err != nil {
			logrus.Fatalf("Error granting permissions: %+v", err)
		}
	}
}

func createUser(conn *pgx.Conn, username string) error {
	rows, err := conn.Query(context.Background(), "SHOW USERS")
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
	if _, err := conn.Exec(context.Background(), "CREATE USER $1", username); err != nil {
		return err
	}

	return nil
}

func createDatabase(conn *pgx.Conn, database string) error {
	rows, err := conn.Query(context.Background(), "SELECT datname FROM pg_database")
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
	if _, err := conn.Exec(context.Background(), fmt.Sprintf("CREATE DATABASE %s", database)); err != nil {
		return err
	}

	return nil
}

func grantPermissions(conn *pgx.Conn, username string, database string) error {
	query := fmt.Sprintf("GRANT ALL ON DATABASE %s TO %s", database, username)
	if _, err := conn.Exec(context.Background(), query); err != nil {
		return err
	}

	logrus.Infof("Granted '%s' permission to read/write to '%s'!", username, database)

	return nil
}
