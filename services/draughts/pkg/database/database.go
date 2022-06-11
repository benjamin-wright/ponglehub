package database

import (
	"context"
	"fmt"
	"strings"
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

func (d *Database) NewGame(player1 string, player2 string) (Game, error) {
	logrus.Infof("Starting new game between %s and %s", player1, player2)

	created := time.Now()

	row := d.conn.QueryRow(
		context.TODO(),
		"INSERT INTO games (player1, player2, turn, created_time, finished) VALUES ($1, $2, 0, $3, false) RETURNING id",
		player1, player2, created,
	)

	game := Game{
		Player1:     uuid.MustParse(player1),
		Player2:     uuid.MustParse(player2),
		Turn:        0,
		CreatedTime: created,
		Finished:    false,
	}

	err := row.Scan(&game.ID)
	if err != nil {
		return Game{}, fmt.Errorf("failed to create database entry: %+v", err)
	}

	return game, nil
}

func (d *Database) NewPieces(pieces []Piece) error {
	query := "INSERT INTO pieces (game, x, y, player, king) VALUES "
	args := []interface{}{}
	index := 1

	for _, piece := range pieces {
		query += fmt.Sprintf("($%d, $%d, $%d, $%d, false),", index, index+1, index+2, index+3)
		args = append(args, piece.Game, piece.X, piece.Y, piece.Player)
		index += 4
	}

	query = strings.TrimSuffix(query, ",")

	cmd, err := d.conn.Exec(context.TODO(), query, args...)
	if err != nil {
		return fmt.Errorf("failed to insert pieces: %+v", err)
	}

	logrus.Infof("Inserted %d pieces", cmd.RowsAffected())

	return nil
}

func (d *Database) LoadGame(id string) (Game, error) {
	row := d.conn.QueryRow(
		context.TODO(),
		"SELECT id, player1, player2, turn, created_time, finished FROM games WHERE id = $1",
		id,
	)

	game := Game{}
	err := row.Scan(&game.ID, &game.Player1, &game.Player2, &game.Turn, &game.CreatedTime, &game.Finished)
	if err != nil {
		return game, fmt.Errorf("failed to load game: %+v", err)
	}

	return game, nil
}

func (d *Database) LoadPieces(game string) ([]Piece, error) {
	rows, err := d.conn.Query(context.TODO(), "SELECT id, game, x, y, player, king FROM pieces WHERE game = $1", game)
	if err != nil {
		return nil, fmt.Errorf("failed to load pieces from database: %+v", err)
	}

	pieces := []Piece{}
	for rows.Next() {
		piece := Piece{}

		err = rows.Scan(&piece.ID, &piece.Game, &piece.X, &piece.Y, &piece.Player, &piece.King)
		if err != nil {
			return nil, fmt.Errorf("failed to parse piece data: %+v", err)
		}

		pieces = append(pieces, piece)
	}

	return pieces, nil
}
