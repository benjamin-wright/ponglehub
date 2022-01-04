package server

import (
	"context"

	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/events/broker/internal/router"
	"ponglehub.co.uk/lib/events"
)

func Start(router *router.Router) (context.CancelFunc, error) {
	cancelFunc, err := events.Listen(80, func(ctx context.Context, event event.Event) {
		logrus.Infof("received event %s from %s: %+v", event.Type(), event.Source(), event)
		urls := router.GetURLs(event.Type())

		for _, url := range urls {
			go func(url string) {
				logrus.Infof("proxying %s, %s -> %s", event.Type(), event.Source(), url)
				client, err := events.New(events.EventsArgs{
					BrokerURL: url,
					Source:    event.Source(),
				})

				if err != nil {
					logrus.Errorf("Failed to create client for %s: %+v", url, err)
					return
				}

				err = client.Proxy(event)
				if err != nil {
					logrus.Errorf("Failed to send event %s %s -> %s: %+v", event.Type(), event.Source(), url, err)
				}
			}(url)
		}
	})

	if err != nil {
		return nil, err
	}

	return cancelFunc, nil
}
