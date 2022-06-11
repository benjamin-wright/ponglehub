package routes

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/games/naughts-and-crosses/pkg/database"
	"ponglehub.co.uk/games/naughts-and-crosses/pkg/rules"
	"ponglehub.co.uk/lib/events"
)

func ListGames(db *database.Database) events.EventRoute {
	return func(userId string, into events.EventParser) ([]events.Response, error) {
		games, err := db.ListGames(userId)
		if err != nil {
			return nil, fmt.Errorf("failed to list games: %+v", err)
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
			return nil, fmt.Errorf("failed to parse payload data from event: %+v", err)
		}

		game, err := db.NewGame(data.Opponent, userId)
		if err != nil {
			return nil, fmt.Errorf("failed to create new game: %+v", err)
		}

		responses := []events.Response{}

		for _, id := range []string{userId, data.Opponent} {
			responses = append(responses, events.Response{
				EventType: "response",
				Data:      map[string]database.Game{"game": game},
				UserId:    id,
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
			return nil, fmt.Errorf("failed to parse payload data from event: %+v", err)
		}

		game, marks, err := db.LoadGame(data.ID)
		if err != nil {
			return []events.Response{{
					EventType: "rejection.response",
					Data:      nil,
					UserId:    userId,
				}},
				fmt.Errorf("failed to load game data: %+v", err)
		}

		return []events.Response{{
				EventType: "response",
				Data:      map[string]interface{}{"game": game, "marks": marks},
				UserId:    userId,
			}},
			nil
	}
}

func Mark(db *database.Database) events.EventRoute {
	return func(userId string, into events.EventParser) ([]events.Response, error) {
		data := struct {
			Game     string `json:"game"`
			Position int    `json:"position"`
		}{}

		err := into(&data)
		if err != nil {
			return nil, fmt.Errorf("failed to parse payload data from event: %+v", err)
		}

		game, marks, err := db.LoadGame(data.Game)
		if err != nil {
			return []events.Response{{
					EventType: "rejection.response",
					Data:      map[string]interface{}{"reason": "server error"},
					UserId:    userId,
				}},
				fmt.Errorf("failed to load game data: %+v", err)
		}

		fail := rules.Validate(game, marks, userId, data.Position)
		if fail != nil {
			return []events.Response{{
					EventType: "rejection.response",
					Data:      map[string]interface{}{"reason": fail.Response()},
					UserId:    userId,
				}},
				errors.New(fail.Log())
		}

		marks = rules.PlaceMark(marks, data.Position, game.Turn)
		if err != nil {
			return []events.Response{{
					EventType: "rejection.response",
					Data:      map[string]interface{}{"reason": "server error"},
					UserId:    userId,
				}},
				fmt.Errorf("failed to load game data: %+v", err)
		}

		winner := rules.IsWinner(marks, data.Position)
		tie := rules.IsTie(marks)
		if winner {
			game.Finished = true
			logrus.Infof("User %d won game %s at position %d", game.Turn, marks, data.Position)
		} else if tie {
			game.Finished = true
			game.Turn = -1
			logrus.Infof("Game %s tied at position %d", marks, data.Position)
		} else {
			logrus.Infof("User %d played mark in game %s at position %d", game.Turn, marks, data.Position)
			game.Turn = rules.NextTurn(game.Turn)
		}

		err = db.SetMarks(data.Game, game.Turn, marks, game.Finished)
		if err != nil {
			return []events.Response{{
					EventType: "rejection.response",
					Data:      map[string]interface{}{"reason": "server error"},
					UserId:    userId,
				}},
				fmt.Errorf("failed to set marks back in database: %+v", err)
		}

		responses := []events.Response{}
		for _, uuid := range []uuid.UUID{game.Player1, game.Player2} {
			responses = append(responses, events.Response{
				EventType: "response",
				Data:      map[string]interface{}{"game": game, "marks": marks},
				UserId:    uuid.String(),
			})
		}

		return responses, nil
	}
}
