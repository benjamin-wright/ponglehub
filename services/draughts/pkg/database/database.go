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
		query += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d),", index, index+1, index+2, index+3, index+4)
		args = append(args, piece.Game, piece.X, piece.Y, piece.Player, piece.King)
		index += 5
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
	defer rows.Close()

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

func (d *Database) Move(game uuid.UUID, piece uuid.UUID, x int16, y int16, king bool, toRemove []uuid.UUID) error {
	tx, err := d.conn.Begin(context.TODO())
	if err != nil {
		return fmt.Errorf("failed to create context: %+v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(context.TODO())
		} else {
			tx.Commit(context.TODO())
		}
	}()

	err = movePiece(tx, game, piece, x, y, king)
	if err != nil {
		return err
	}

	err = removePieces(tx, game, toRemove)
	if err != nil {
		return err
	}

	return nil
}

func movePiece(tx pgx.Tx, game uuid.UUID, piece uuid.UUID, x int16, y int16, king bool) error {
	kingQuery := ""

	if king {
		kingQuery = ", king = true"
	}

	_, err := tx.Exec(context.TODO(), fmt.Sprintf("UPDATE pieces SET x = $1, y = $2 %s WHERE id = $3", kingQuery), x, y, piece)
	if err != nil {
		return fmt.Errorf("failed to update piece: %+v", err)
	}

	return nil
}

func removePieces(tx pgx.Tx, game uuid.UUID, ids []uuid.UUID) error {
	placeholders := make([]string, len(ids))
	params := make([]interface{}, len(ids)+1)
	params[0] = game

	for idx, id := range ids {
		placeholders[idx] = fmt.Sprintf("$%d", idx+2)
		params[idx+1] = id
	}

	query := "DELETE FROM pieces WHERE game = $1 AND id IN (\"" + strings.Join(placeholders, ", ") + "\")"

	_, err := tx.Exec(context.TODO(), query, params...)
	if err != nil {
		return fmt.Errorf("failed to delete pieces: %+v", err)
	}

	return nil
}
