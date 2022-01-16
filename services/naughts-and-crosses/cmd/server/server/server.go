package server

import (
	"context"
	"errors"
	"fmt"

	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/games/naughts-and-crosses/pkg/database"
	"ponglehub.co.uk/lib/events"
)

type Server struct {
	cancelFunc context.CancelFunc
}

func New(client *events.Events, db *database.Database) (*Server, error) {
	cancelFunc, err := events.Listen(80, func(ctx context.Context, event event.Event) {
		var err error

		switch event.Type() {
		case "naughts-and-crosses.new-game":
			err = newGame(client, db, event)
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
	Player1 string `json:"player1"`
	Player2 string `json:"player2"`
}

func newGame(client *events.Events, db *database.Database, event event.Event) error {
	data := NewGameEvent{}
	event.DataAs(&data)

	id, err := db.NewGame(data.Player1, data.Player2)
	if err != nil {
		return fmt.Errorf("failed to create new game: %+v", err)
	}

	err = client.Send("naughts-and-crosses.new-game-id", map[string]string{"id": id})
	if err != nil {
		return fmt.Errorf("failed to send new game id event: %+v", err)
	}

	return nil
}
