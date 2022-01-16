package server

import (
	"context"

	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/lib/events"
)

func Start() (context.CancelFunc, error) {
	cancelFunc, err := events.Listen(80, func(ctx context.Context, event event.Event) {
		logrus.Infof("received event %s from %s", event.Type(), event.Source())
	})

	if err != nil {
		return nil, err
	}

	return cancelFunc, nil
}
