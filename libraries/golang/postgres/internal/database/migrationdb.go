package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"

	"ponglehub.co.uk/auth/postgres/pkg/connect"
)

type MigrationConn struct {
	conn *pgx.Conn
}

func NewMigrationConn(cfg connect.ConnectConfig) (*MigrationConn, error) {
	conn, err := connect.Connect(cfg)
	if err != nil {
		return nil, err
	}

	return &MigrationConn{conn}, nil
}

func (d *MigrationConn) Stop() {
	d.conn.Close(context.Background())
}

func (d *MigrationConn) EnsureMigrationTable() error {
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

func (d *MigrationConn) HasMigration(id int) bool {
	var found int
	err := d.conn.QueryRow(context.Background(), "SELECT id FROM migrations WHERE id == $1", id).Scan(&found)
	return err == nil
}

func (d *MigrationConn) AddMigration(id int) error {
	_, err := d.conn.Exec(context.Background(), "INSERT INTO migrations (id) VALUES ($1)", id)
	return err
}

func (d *MigrationConn) RunMigration(query string) error {
	_, err := d.conn.Exec(context.TODO(), query)

	return err
}

func (d *MigrationConn) GetTables() ([]string, error) {
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

func (d *MigrationConn) GetTableSchema(tableName string) (map[string]string, error) {
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

func (d *MigrationConn) GetContents(tableName string) ([][]interface{}, error) {
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
