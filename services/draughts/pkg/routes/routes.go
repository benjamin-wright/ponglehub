package routes

import (
	"fmt"

	"github.com/google/uuid"
	"ponglehub.co.uk/games/draughts/pkg/database"
	"ponglehub.co.uk/games/draughts/pkg/rules"
	"ponglehub.co.uk/lib/events"
)

func ListGames(db *database.Database) events.EventRoute {
	return func(userId string, into events.EventParser) ([]events.Response, error) {
		games, err := db.ListGames(userId)
		if err != nil {
			return nil, fmt.Errorf("error listing games: %+v", err)
		}

		return []events.Response{{
			EventType: "response",
			Data:      map[string]interface{}{"games": games},
			UserId:    userId,
		}}, nil
	}
}

func NewGame(db *database.Database) events.EventRoute {
	return func(userId string, into events.EventParser) ([]events.Response, error) {
		data := struct {
			Opponent string `json:"opponent"`
		}{}

		err := into(&data)
		if err != nil {
			return nil, fmt.Errorf("failed to parse new game event data: %+v", err)
		}

		game, err := db.NewGame(userId, data.Opponent)
		if err != nil {
			return nil, fmt.Errorf("failed to create new game: %+v", err)
		}

		pieces := rules.NewGame(game.ID)
		err = db.NewPieces(pieces)
		if err != nil {
			return nil, fmt.Errorf("failed to create new pieces: %+v", err)
		}

		responses := []events.Response{}
		for _, id := range []string{userId, data.Opponent} {
			responses = append(responses, events.Response{
				EventType: "response",
				Data: map[string]interface{}{
					"game": game,
				},
				UserId: id,
			})
		}

		return responses, nil
	}
}

func LoadGame(db *database.Database) events.EventRoute {
	return func(userId string, into events.EventParser) ([]events.Response, error) {
		data := struct {
			ID string `json:"id"`
		}{}
		err := into(&data)
		if err != nil {
			return nil, fmt.Errorf("failed to parse load game args %s: %+v", data.ID, err)
		}

		game, err := db.LoadGame(data.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to load game data %s: %+v", data.ID, err)
		}

		pieces, err := db.LoadPieces(data.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to load game pieces %s: %+v", data.ID, err)
		}

		return []events.Response{{
			EventType: "response",
			Data: map[string]interface{}{
				"game":   game,
				"pieces": pieces,
			},
			UserId: userId,
		}}, nil
	}
}

func Move(db *database.Database) events.EventRoute {
	return func(userId string, into events.EventParser) ([]events.Response, error) {
		data := struct {
			Game  string       `json:"game"`
			Moves []rules.Move `json:"moves"`
		}{}

		err := into(&data)
		if err != nil {
			return nil, fmt.Errorf("failed to parse move data: %+v", err)
		}

		game, err := db.LoadGame(data.Game)
		if err != nil {
			return nil, fmt.Errorf("failed to load game: %+v", err)
		}

		pieces, err := db.LoadPieces(data.Game)
		if err != nil {
			return nil, fmt.Errorf("failed to load pieces: %+v", err)
		}

		if !rules.IsYourTurn(userId, game) {
			return []events.Response{{
				EventType: "rejection.response",
				Data:      map[string]interface{}{"message": "it's not your turn"},
				UserId:    userId,
			}}, fmt.Errorf("user %s made a move when it wasn't their turn", userId)
		}

		result, err := rules.Process(data.Moves, pieces)
		if err != nil {
			return []events.Response{{
				EventType: "rejection.response",
				Data:      map[string]interface{}{"message": "invalid move"},
				UserId:    userId,
			}}, fmt.Errorf("user %s made an invalid move: %+v", userId, err)
		}

		err = db.Move(game.ID, result.Piece, result.NewX, result.NewY, result.King, result.ToRemove)
		if err != nil {
			return []events.Response{{
				EventType: "rejection.response",
				Data:      map[string]interface{}{"message": "server error"},
				UserId:    userId,
			}}, fmt.Errorf("failed to process user %s move: %+v", userId, err)
		}

		pieces, err = db.LoadPieces(game.ID.String())
		if err != nil {
			return []events.Response{{
				EventType: "rejection.response",
				Data:      map[string]interface{}{"message": "server error"},
				UserId:    userId,
			}}, fmt.Errorf("failed to fetch pieces after user %s move: %+v", userId, err)
		}

		responses := []events.Response{}

		for _, id := range []uuid.UUID{game.Player1, game.Player2} {
			responses = append(responses, events.Response{
				EventType: "response",
				Data:      map[string]interface{}{"pieces": pieces},
				UserId:    id.String(),
			})
		}

		return responses, nil
	}
}
