package server

import (
	"context"
	"errors"
	"fmt"

	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/games/draughts/pkg/database"
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
		case "draughts.list-games":
			err = listGames(userId, client, db)
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

func listGames(userId string, client *events.Events, db *database.Database) error {
	games, err := db.ListGames(userId)
	if err != nil {
		return fmt.Errorf("error listing games: %+v", err)
	}

	err = client.Send(
		"draughts.list-games.response",
		map[string]interface{}{"games": games},
		map[string]interface{}{"userid": userId},
	)
	if err != nil {
		return fmt.Errorf("failed to send list games response event: %+v", err)
	}

	return nil
}
