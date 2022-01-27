package server

import (
	"context"
	"errors"
	"fmt"

	"github.com/cloudevents/sdk-go/v2/event"
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

		switch event.Type() {
		case "naughts-and-crosses.new-game":
			err = newGame(client, db, userId, event)
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

type NewGameEvent struct {
	Opponent string `json:"opponent"`
}

func newGame(client *events.Events, db *database.Database, userId string, event event.Event) error {
	data := NewGameEvent{}
	err := event.DataAs(&data)
	if err != nil {
		return fmt.Errorf("failed to parse payload data from event: %+v", err)
	}

	id, err := db.NewGame(data.Opponent, userId)
	if err != nil {
		return fmt.Errorf("failed to create new game: %+v", err)
	}

	err = client.Send("naughts-and-crosses.new-game-id", map[string]string{"id": id})
	if err != nil {
		return fmt.Errorf("failed to send new game id event: %+v", err)
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

	game, err := db.GetMarks(data.Game)
	if err != nil {
		return fmt.Errorf("failed to fetch game marks: %+v", err)
	}

	err = rules.ValidateMark(game, userId, data.Position)
	if err != nil {
		return fmt.Errorf("failed to fetch game marks: %+v", err)
	}

	err = db.SetMark(data.Game, data.Position)
	if err != nil {
		return fmt.Errorf("failed to set new mark: %+v", err)
	}

	err = db.ChangeTurn(game)
	if err != nil {
		return fmt.Errorf("failed to change turn: %+v", err)
	}

	game.SetMark(data.Position)

	err = client.Send("naughts-and-crosses.new-marks", map[string]interface{}{
		"game":    data.Game,
		"turn":    userId,
		"player1": game.Player1Marks,
		"player2": game.Player2Marks,
	})
	if err != nil {
		return fmt.Errorf("failed to send new game id event: %+v", err)
	}

	return nil
}
