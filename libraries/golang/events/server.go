package events

import (
	"context"
	"fmt"

	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/sirupsen/logrus"
)

type Response struct {
	EventType string
	Data      interface{}
	UserId    string
}

type EventParser func(obj interface{}) error

type EventRoutes map[string]EventRoute

type EventRoute func(userId string, into EventParser) ([]Response, error)

type ServeParams struct {
	BrokerEnv string
	BrokerURL string
	Source    string
	Routes    map[string]EventRoute
}

func Serve(params ServeParams) (context.CancelFunc, error) {
	client, err := New(EventsArgs{
		BrokerEnv: params.BrokerEnv,
		BrokerURL: params.BrokerURL,
		Source:    params.Source,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create client connection: %+v", err)
	}

	cancelFunc, err := Listen(80, func(ctx context.Context, event event.Event) {
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

		route, ok := params.Routes[event.Type()]
		if !ok {
			logrus.Errorf("unexpected event type: %s", event.Type())
			return
		}

		responses, err := route(userId, event.DataAs)
		if err != nil {
			logrus.Errorf("error processing event %s: %+v", event.Type(), err)
		}

		for _, response := range responses {
			err = client.Send(
				fmt.Sprintf("%s.%s", event.Type(), response.EventType),
				response.Data,
				map[string]interface{}{"userid": response.UserId},
			)

			if err != nil {
				logrus.Errorf("failed to send \"%s\" response to event \"%s\": %+v", response.EventType, event.Type(), err)
			}
		}
	})

	if err != nil {
		return nil, err
	}

	return cancelFunc, nil
}
