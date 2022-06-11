package routes

import (
	"fmt"

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
