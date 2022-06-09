package database

import (
	"context"
	"fmt"
	"time"

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
	logrus.Info("Clearing all game data")
	cmd, err := d.conn.Exec(context.TODO(), "DELETE FROM games")
	if err != nil {
		return fmt.Errorf("error clearing games table: %+v", err)
	}
	logrus.Infof("Cleared %d rows from 'games'", cmd.RowsAffected())

	cmd, err = d.conn.Exec(context.TODO(), "DELETE FROM pieces")
	if err != nil {
		return fmt.Errorf("error clearing pieces table: %+v", err)
	}
	logrus.Infof("Cleared %d rows from 'pieces'", cmd.RowsAffected())

	return nil
}

type Game struct {
	ID          uuid.UUID `json:"id"`
	Player1     uuid.UUID `json:"player1"`
	Player2     uuid.UUID `json:"player2"`
	Turn        int16     `json:"turn"`
	CreatedTime time.Time `json:"createdTime"`
	Finished    bool      `json:"finished"`
}

func (d *Database) ListGames(user string) ([]Game, error) {
	logrus.Infof("Listing games for user %s", user)
	rows, err := d.conn.Query(context.TODO(), "SELECT id, player1, player2, turn, created_time, finished FROM games WHERE player1=$1 OR player2=$1", user)
	if err != nil {
		return nil, fmt.Errorf("error fetching games data: %+v", err)
	}
	defer rows.Close()

	games := []Game{}

	for rows.Next() {
		game := Game{}
		err = rows.Scan(&game.ID, &game.Player1, &game.Player2, &game.Turn, &game.CreatedTime, &game.Finished)
		if err != nil {
			return nil, fmt.Errorf("error parsing game data: %+v", err)
		}

		games = append(games, game)
	}

	return games, nil
}

func (d *Database) InsertGame(game Game) error {
	logrus.Warnf("Inserting test game data: %+v", game)
	_, err := d.conn.Exec(
		context.TODO(),
		"INSERT INTO games (id, player1, player2, turn, created_time, finished) VALUES ($1, $2, $3, $4, $5, $6)",
		game.ID, game.Player1, game.Player2, game.Turn, game.CreatedTime, game.Finished,
	)
	if err != nil {
		return fmt.Errorf("failed to insert new game: %+v", err)
	}

	return nil
}
