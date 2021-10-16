package database

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"path"

	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
)

type Database struct {
	conn  *pgx.Conn
	admin bool
}

func getTlsConfig(certsDir string, username string) (*tls.Config, error) {
	clientCrt := path.Join(certsDir, fmt.Sprintf("client.%s.crt", username))
	clientKey := path.Join(certsDir, fmt.Sprintf("client.%s.key", username))
	caCrt := path.Join(certsDir, "ca.crt")

	cert, err := tls.LoadX509KeyPair(clientCrt, clientKey)
	if err != nil {
		return nil, fmt.Errorf("failed to load client certificate: %v", err)
	}

	CACert, err := ioutil.ReadFile(caCrt)
	if err != nil {
		return nil, fmt.Errorf("failed to load server certificate: %v", err)
	}

	CACertPool := x509.NewCertPool()
	CACertPool.AppendCertsFromPEM(CACert)

	return &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            CACertPool,
		InsecureSkipVerify: true,
	}, nil
}

func New(host string, port int, username string, password string, database string, certsDir string) (*Database, error) {
	pgxConfig, err := pgx.ParseConfig(fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", username, password, host, port, database))
	if err != nil {
		return nil, err
	}

	tlsConfig, err := getTlsConfig(certsDir, username)
	if err != nil {
		return nil, err
	}

	pgxConfig.TLSConfig = tlsConfig
	conn := getConnection(pgxConfig)
	if conn == nil {
		return nil, errors.New("failed to create connection, exiting")
	}

	return &Database{
		conn:  conn,
		admin: false,
	}, nil
}

func Admin(host string, port int, certsDir string) (*Database, error) {
	pgxConfig, err := pgx.ParseConfig(fmt.Sprintf("postgresql://%s@%s:%d", "root", host, port))
	if err != nil {
		return nil, err
	}

	tlsConfig, err := getTlsConfig(certsDir, "root")
	if err != nil {
		return nil, err
	}

	pgxConfig.TLSConfig = tlsConfig
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

func (d *Database) CreateUser(username string, password string) error {
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
	if _, err := d.conn.Exec(context.Background(), "CREATE USER $1 WITH PASSWORD $2", username, password); err != nil {
		return err
	}

	return nil
}

func (d *Database) DropUser(username string) error {
	if !d.admin {
		return errors.New("cannot call DropUser on non-admin connection")
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
			rows.Close()

			logrus.Infof("Deleting user %s!", username)
			_, err := d.conn.Exec(context.Background(), "DROP USER $1", username)
			return err
		}
	}

	logrus.Infof("User %s doesn't exist!", username)
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

func (d *Database) DropDatabase(database string) error {
	if !d.admin {
		return errors.New("cannot call DropDatabase on non-admin connection")
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
			rows.Close()

			logrus.Infof("Dropping database %s!", database)
			_, err := d.conn.Exec(context.Background(), fmt.Sprintf("DROP DATABASE %s", database))
			return err
		}
	}

	logrus.Infof("Database %s didn't exist!", database)
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

func (d *Database) RevokePermissions(username string, database string) error {
	if !d.admin {
		return errors.New("cannot call RevokePermissions on non-admin connection")
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
			rows.Close()

			query := fmt.Sprintf("REVOKE ALL ON DATABASE %s FROM %s", database, username)
			if _, err := d.conn.Exec(context.Background(), query); err != nil {
				return err
			}

			logrus.Infof("Revoked '%s' permission to read/write from '%s'!", username, database)
			return nil
		}
	}

	logrus.Infof("User '%s' doesn't exist!", username)
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

func (d *Database) HasMigration(id int) bool {
	var found int
	err := d.conn.QueryRow(context.Background(), "SELECT id FROM migrations WHERE id == $1", id).Scan(&found)
	return err == nil
}

func (d *Database) AddMigration(id int) error {
	_, err := d.conn.Exec(context.Background(), "INSERT INTO migrations (id) VALUES ($1)", id)
	return err
}

func (d *Database) RunMigration(query string) error {
	_, err := d.conn.Exec(context.TODO(), query)

	return err
}

func (d *Database) GetTables() ([]string, error) {
	rows, err := d.conn.Query(context.TODO(), "SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname != 'pg_catalog' AND schemaname != 'information_schema'")
	if err != nil {
		return nil, err
	}

	names := []string{}

	for rows.Next() {
		name := ""
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}

		names = append(names, name)
	}

	return names, nil
}

func (d *Database) GetTableSchema(tableName string) (map[string]string, error) {
	rows, err := d.conn.Query(context.TODO(), "SELECT column_name, data_type FROM information_schema.columns WHERE table_name = $1", tableName)
	if err != nil {
		return nil, err
	}

	columns := map[string]string{}

	for rows.Next() {
		var column string
		var dataType string

		if err = rows.Scan(&column, &dataType); err != nil {
			return nil, err
		}

		columns[column] = dataType
	}

	return columns, err
}

func (d *Database) GetContents(tableName string) ([][]interface{}, error) {
	rows, err := d.conn.Query(context.TODO(), fmt.Sprintf("SELECT * FROM %s", pgx.Identifier{tableName}.Sanitize()))
	if err != nil {
		return nil, err
	}

	contents := [][]interface{}{}

	for rows.Next() {
		row, err := rows.Values()
		if err != nil {
			return nil, err
		}

		contents = append(contents, row)
	}

	return contents, err
}
