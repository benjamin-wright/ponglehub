package server

import (
	"context"
	"errors"
	"fmt"

	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/games/naughts-and-crosses/pkg/database"
	"ponglehub.co.uk/games/naughts-and-crosses/pkg/rules"
	"ponglehub.co.uk/lib/events"
)

type Server struct {
	cancelFunc context.CancelFunc
}

func New(client *events.Events, db *database.Database) (*Server, error) {
	cancelFunc, err := events.Listen(80, func(ctx context.Context, event event.Event) {
		var err error

		userIdObj, err := event.Context.GetExtension("userid")
		if err != nil {
			logrus.Errorf("failed to get user id from event: %+v", err)
			return
		}

		userId, ok := userIdObj.(string)
		if !ok {
			logrus.Errorf("expected user id to be a string, got %T", userId)
			return
		}

		logrus.Infof("Got event: %s", event.Type())

		switch event.Type() {
		case "naughts-and-crosses.list-games":
			err = listGames(client, db, userId)
		case "naughts-and-crosses.new-game":
			err = newGame(client, db, userId, event)
		case "naughts-and-crosses.load-game":
			err = loadGame(client, db, userId, event)
		case "naughts-and-crosses.mark":
			err = mark(client, db, userId, event)
		default:
			err = errors.New("unrecognised event type")
		}

		if err != nil {
			logrus.Errorf("Failed to process event type %s: %+v", event.Type(), err)
		}
	})

	if err != nil {
		return nil, err
	}

	return &Server{
		cancelFunc,
	}, nil
}

func (s *Server) Stop() {
	s.cancelFunc()
}

func listGames(client *events.Events, db *database.Database, userId string) error {
	games, err := db.ListGames(userId)
	if err != nil {
		return fmt.Errorf("failed to list games: %+v", err)
	}

	err = client.Send(
		"naughts-and-crosses.list-games.response",
		map[string]interface{}{"games": games},
		map[string]interface{}{"userid": userId},
	)
	if err != nil {
		return fmt.Errorf("failed to send new game id event: %+v", err)
	}

	return nil
}

type NewGameEvent struct {
	Opponent string `json:"opponent"`
}

func newGame(client *events.Events, db *database.Database, userId string, event event.Event) error {
	data := NewGameEvent{}
	err := event.DataAs(&data)
	if err != nil {
		return fmt.Errorf("failed to parse payload data from event: %+v", err)
	}

	game, err := db.NewGame(data.Opponent, userId)
	if err != nil {
		return fmt.Errorf("failed to create new game: %+v", err)
	}

	for _, id := range []string{userId, data.Opponent} {
		err = client.Send(
			"naughts-and-crosses.new-game.response",
			map[string]database.Game{"game": game},
			map[string]interface{}{"userid": id},
		)
		if err != nil {
			return fmt.Errorf("failed to send new game response for %s: %+v", id, err)
		}
	}

	return nil
}

type LoadGameEvent struct {
	ID string `json:"id"`
}

func loadGame(client *events.Events, db *database.Database, userId string, event event.Event) error {
	data := LoadGameEvent{}
	err := event.DataAs(&data)
	if err != nil {
		return fmt.Errorf("failed to parse payload data from event: %+v", err)
	}

	game, marks, err := db.LoadGame(data.ID)
	if err != nil {
		client.Send(
			"naughts-and-crosses.load-game.rejection.response",
			nil,
			map[string]interface{}{"userid": userId},
		)
		return fmt.Errorf("failed to load game data: %+v", err)
	}

	err = client.Send(
		"naughts-and-crosses.load-game.response",
		map[string]interface{}{"game": game, "marks": marks},
		map[string]interface{}{"userid": userId},
	)
	if err != nil {
		return fmt.Errorf("failed to send load game response for %s: %+v", userId, err)
	}

	return nil
}

type MarkEvent struct {
	Game     string `json:"game"`
	Position int    `json:"position"`
}

func mark(client *events.Events, db *database.Database, userId string, event event.Event) error {
	data := MarkEvent{}
	err := event.DataAs(&data)
	if err != nil {
		return fmt.Errorf("failed to parse payload data from event: %+v", err)
	}

	game, marks, err := db.LoadGame(data.Game)
	if err != nil {
		client.Send(
			"naughts-and-crosses.mark.rejection.response",
			map[string]interface{}{"reason": "server error"},
			map[string]interface{}{"userid": userId},
		)
		return fmt.Errorf("failed to load game data: %+v", err)
	}

	fail := rules.Validate(game, marks, userId, data.Position)
	if fail != nil {
		client.Send(
			"naughts-and-crosses.mark.rejection.response",
			map[string]interface{}{"reason": fail.Response()},
			map[string]interface{}{"userid": userId},
		)
		return errors.New(fail.Log())
	}

	marks = rules.PlaceMark(marks, data.Position, game.Turn)
	if err != nil {
		client.Send(
			"naughts-and-crosses.mark.rejection.response",
			map[string]interface{}{"reason": "server error"},
			map[string]interface{}{"userid": userId},
		)
		return fmt.Errorf("failed to load game data: %+v", err)
	}

	game.Finished = rules.IsWinner(marks, data.Position)
	if game.Finished {
		logrus.Infof("User %d won game %s at position %d", game.Turn, marks, data.Position)
	} else {
		logrus.Infof("User %d played mark in game %s at position %d", game.Turn, marks, data.Position)
		game.Turn = rules.NextTurn(game.Turn)
	}

	err = db.SetMarks(data.Game, game.Turn, marks, game.Finished)
	if err != nil {
		client.Send(
			"naughts-and-crosses.mark.rejection.response",
			map[string]interface{}{"reason": "server error"},
			map[string]interface{}{"userid": userId},
		)
		return fmt.Errorf("failed to set marks back in database: %+v", err)
	}

	for _, uuid := range []uuid.UUID{game.Player1, game.Player2} {
		err = client.Send(
			"naughts-and-crosses.mark.response",
			map[string]interface{}{"game": game, "marks": marks},
			map[string]interface{}{"userid": uuid.String()},
		)
		if err != nil {
			logrus.Errorf("failed to send mark response for %s - %s: %+v", userId, data.Game, err)
		}
	}

	return nil
}
