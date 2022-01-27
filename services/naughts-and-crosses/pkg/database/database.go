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

	cmd, err = d.conn.Exec(context.TODO(), "DELETE FROM marks")
	if err != nil {
		return fmt.Errorf("error clearing games table: %+v", err)
	}
	logrus.Infof("Cleared %d rows from 'marks'", cmd.RowsAffected())

	return nil
}

func (d *Database) NewGame(player1 string, player2 string) (string, error) {
	row := d.conn.QueryRow(context.TODO(), "INSERT INTO games (player1, player2, turn) VALUES ($1, $2, 0) RETURNING id;", player1, player2)

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

type GameState struct {
	ID           uuid.UUID
	Player1      uuid.UUID
	Player2      uuid.UUID
	Turn         int
	Player1Marks []int
	Player2Marks []int
}

func (g *GameState) CurrentPlayer() string {
	switch g.Turn {
	case 0:
		return g.Player1.String()
	case 1:
		return g.Player2.String()
	}

	return ""
}

func (g *GameState) SetMark(position int) {
	switch g.Turn {
	case 0:
		g.Player1Marks = append(g.Player1Marks, position)
	case 1:
		g.Player2Marks = append(g.Player2Marks, position)
	}
}

func (d *Database) GetMarks(game string) (*GameState, error) {
	state := GameState{
		ID: uuid.MustParse(game),
	}

	row := d.conn.QueryRow(context.TODO(), "SELECT player1, player2, turn FROM games WHERE id=$1", game)
	if err := row.Scan(&state.Player1, &state.Player2, &state.Turn); err != nil {
		return nil, fmt.Errorf("error fetching data for game '%s': %+v", game, err)
	}

	rows, err := d.conn.Query(context.TODO(), "SELECT player, position FROM marks WHERE game_id=$1", game)
	if err != nil {
		return nil, fmt.Errorf("error fetching marks for game '%s': %+v", game, err)
	}

	for rows.Next() {
		player := uuid.UUID{}
		position := 0

		err = rows.Scan(&player, &position)
		if err != nil {
			return nil, fmt.Errorf("error parsing marks data for game '%s': %+v", game, err)
		}

		switch player.String() {
		case state.Player1.String():
			state.Player1Marks = append(state.Player1Marks, position)
		case state.Player2.String():
			state.Player2Marks = append(state.Player2Marks, position)
		}
	}

	return &state, nil
}

func (d *Database) SetMark(game string, turn int) error {
	cmd, err := d.conn.Exec(context.TODO(), "update games SET turn=$1 WHERE id=$2", turn, game)
	if err != nil {
		return fmt.Errorf("error changing turn: %+v", err)
	}

	if cmd.RowsAffected() != 1 {
		return fmt.Errorf("'set turn' for game %s should have affected 1 row, got %d", game, cmd.RowsAffected())
	}

	return nil
}

func (d *Database) ChangeTurn(game *GameState) error {
	newTurn := 1
	if game.Turn == 1 {
		newTurn = 0
	}

	cmd, err := d.conn.Exec(context.TODO(), "update games SET turn=$1 WHERE id=$2", newTurn, game.ID)
	if err != nil {
		return fmt.Errorf("error changing turn: %+v", err)
	}

	if cmd.RowsAffected() != 1 {
		return fmt.Errorf("'set turn' for game %s should have affected 1 row, got %d", game.ID, cmd.RowsAffected())
	}

	return nil
}
