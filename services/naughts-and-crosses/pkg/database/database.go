package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/pgtype"
	"github.com/jackc/pgx/v4"
	"ponglehub.co.uk/lib/postgres/pkg/connect"
)

type Database struct {
	conn *pgx.Conn
}

func New() (*Database, error) {
	cfg, err := connect.ConfigFromEnv()
	if err != nil {
		return nil, fmt.Errorf("failed to load config from environment: %+v", err)
	}

	conn, err := connect.Connect(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %+v", err)
	}

	return &Database{conn}, nil
}

func (d *Database) NewGame(player1 string, player2 string) (string, error) {
	row := d.conn.QueryRow(context.TODO(), "INSERT INTO games (player1, player2) VALUES ($1, $2) RETURNING id;", player1, player2)

	id := map[string]pgtype.UUID{"id": pgtype.UUID{}}

	err := row.Scan(&id)
	if err != nil {
		return "", fmt.Errorf("failed to add new game: %+v", err)
	}

	return string(id["id"].), nil
}
