package database

import (
	"context"
	"errors"
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

	return nil
}

func (d *Database) NewGame(player1 string, player2 string) (Game, error) {
	logrus.Infof("Creating new game for %s vs %s", player1, player2)
	created := time.Now()

	row := d.conn.QueryRow(
		context.TODO(),
		"INSERT INTO games (player1, player2, created_time, turn, marks) VALUES ($1, $2, $3, 0, '---------') RETURNING id;",
		player1,
		player2,
		created,
	)

	var id uuid.UUID

	err := row.Scan(&id)
	if err != nil {
		return Game{}, fmt.Errorf("failed to add new game: %+v", err)
	}

	return Game{
		ID:      id,
		Player1: uuid.MustParse(player1),
		Player2: uuid.MustParse(player2),
		Turn:    0,
		Created: created,
	}, nil
}

func (d *Database) InsertGame(game Game, marks string) error {
	logrus.Infof("Inserting game for %s vs %s", game.Player1, game.Player2)
	_, err := d.conn.Exec(
		context.TODO(),
		"INSERT INTO games (id, player1, player2, created_time, turn, marks) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;",
		game.ID, game.Player1, game.Player2, game.Created, game.Turn, marks,
	)
	if err != nil {
		return fmt.Errorf("failed to add new game: %+v", err)
	}

	return nil
}

type Game struct {
	ID      uuid.UUID
	Player1 uuid.UUID
	Player2 uuid.UUID
	Created time.Time
	Turn    int16
}

func (d *Database) ListGames(player string) ([]Game, error) {
	logrus.Infof("Listing games for user %s", player)
	rows, err := d.conn.Query(context.TODO(), "SELECT id, player1, player2, created_time, turn FROM games WHERE player1=$1 OR player2=$1", player)
	if err != nil {
		return nil, fmt.Errorf("error fetching games data: %+v", err)
	}
	defer rows.Close()

	games := []Game{}

	for rows.Next() {
		game := Game{}
		err := rows.Scan(&game.ID, &game.Player1, &game.Player2, &game.Created, &game.Turn)
		if err != nil {
			return nil, fmt.Errorf("error parsing game data: %+v", err)
		}

		games = append(games, game)
	}

	return games, nil
}

func (d *Database) LoadGame(id string) (*Game, string, error) {
	logrus.Infof("Loading game %s", id)
	rows, err := d.conn.Query(context.TODO(), "SELECT id, player1, player2, created_time, turn, marks FROM games WHERE id=$1", id)
	if err != nil {
		return nil, "", fmt.Errorf("error fetching game data: %+v", err)
	}
	defer rows.Close()

	game := Game{}
	var marks string

	if !rows.Next() {
		return nil, "", errors.New("game not found")
	}

	err = rows.Scan(&game.ID, &game.Player1, &game.Player2, &game.Created, &game.Turn, &marks)
	if err != nil {
		return nil, "", fmt.Errorf("error parsing game data: %+v", err)
	}

	return &game, marks, nil
}

func (d *Database) SetMarks(id string, turn int16, marks string) error {
	logrus.Infof("Updating game %s", id)

	_, err := d.conn.Exec(context.TODO(), "UPDATE games SET turn=$1, marks=$2 WHERE id=$3", turn, marks, id)
	if err != nil {
		return fmt.Errorf("error setting mark data: %+v", err)
	}

	return nil
}
