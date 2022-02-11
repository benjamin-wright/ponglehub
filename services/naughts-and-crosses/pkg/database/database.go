package database

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
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

func (d *Database) Clear() error {
	cmd, err := d.conn.Exec(context.TODO(), "DELETE FROM games")
	if err != nil {
		return fmt.Errorf("error clearing games table: %+v", err)
	}
	logrus.Infof("Cleared %d rows from 'games'", cmd.RowsAffected())

	return nil
}

func (d *Database) NewGame(player1 string, player2 string) (string, error) {
	row := d.conn.QueryRow(context.TODO(), "INSERT INTO games (player1, player2, turn, marks) VALUES ($1, $2, 0, '---------') RETURNING id;", player1, player2)

	var id uuid.UUID

	err := row.Scan(&id)
	if err != nil {
		return "", fmt.Errorf("failed to add new game: %+v", err)
	}

	return id.String(), nil
}

type Game struct {
	ID      uuid.UUID
	Player1 uuid.UUID
	Player2 uuid.UUID
	Turn    int16
}

func (d *Database) ListGames(player string) ([]Game, error) {
	rows, err := d.conn.Query(context.TODO(), "SELECT id, player1, player2, turn FROM games WHERE player1=$1 OR player2=$1", player)
	if err != nil {
		return nil, fmt.Errorf("error fetching games data: %+v", err)
	}

	games := []Game{}

	for rows.Next() {
		game := Game{}
		err := rows.Scan(&game.ID, &game.Player1, &game.Player2, &game.Turn)
		if err != nil {
			return nil, fmt.Errorf("error parsing game data: %+v", err)
		}

		games = append(games, game)
	}

	return games, nil
}
