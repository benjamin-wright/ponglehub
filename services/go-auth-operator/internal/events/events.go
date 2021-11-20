package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/cloudevents/sdk-go/v2/protocol/http"
	"github.com/sirupsen/logrus"
)

type Events struct {
	ctx    context.Context
	sender cloudevents.Client
}

type User struct {
	Name            string `json:"name"`
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ID              string `json:"id"`
	Pending         bool   `json:"pending"`
	ResourceVersion string `json:"resource_version"`
}

func (a User) Equals(user User) bool {
	return a.Email == user.Email &&
		a.Username == user.Username &&
		a.Password == user.Password
}

type UserEvent struct {
	Type string
	User User
}

type UserEventHandler func(event UserEvent)

func New(brokerUrl string) (*Events, error) {
	ctx := cloudevents.ContextWithTarget(context.Background(), brokerUrl)

	p, err := cloudevents.NewHTTP()
	if err != nil {
		return nil, fmt.Errorf("error connecting to NATS: %+v", err)
	}

	p.Client.Timeout = time.Second

	client, err := cloudevents.NewClient(p, cloudevents.WithTimeNow(), cloudevents.WithUUIDs())
	if err != nil {
		return nil, fmt.Errorf("error creating cloudevents instance: %+v", err)
	}

	return &Events{
		ctx:    ctx,
		sender: client,
	}, nil
}

func Listen(handler UserEventHandler) (context.CancelFunc, error) {
	p, err := cloudevents.NewHTTP(cloudevents.WithPort(80))
	if err != nil {
		return nil, fmt.Errorf("failed to create protocol: %s", err.Error())
	}

	client, err := cloudevents.NewClient(p)
	if err != nil {
		return nil, fmt.Errorf("failed to create client, %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		err := client.StartReceiver(ctx, func(ctx context.Context, e event.Event) {
			logrus.Infof("Received event: %s", e.Type())

			var user User
			err = json.Unmarshal(e.Data(), &user)
			if err != nil {
				logrus.Errorf("Error parsing user set data: %+v", err)
				return
			}

			handler(UserEvent{Type: e.Type(), User: user})
		})

		if err != nil && ctx.Err() == nil {
			logrus.Fatalf("Error in event listener: %+v", err)
		} else {
			logrus.Infof("Stopped event listener")
		}
	}()

	return cancel, nil
}

func (e *Events) sendEvent(sender cloudevents.Client, eventType string, data interface{}) error {
	event := cloudevents.NewEvent()
	event.SetType(eventType)
	event.SetSource("auth-operator")
	err := event.SetData(cloudevents.ApplicationJSON, data)
	if err != nil {
		return fmt.Errorf("failed to serialize event data: %+v", err)
	}

	ctx := cloudevents.ContextWithRetriesConstantBackoff(e.ctx, time.Second, 20)
	res := sender.Send(ctx, event)

	if cloudevents.IsUndelivered(res) {
		return fmt.Errorf("failed to send event: %v", res.Error())
	}

	if !cloudevents.IsACK(res) {
		return fmt.Errorf("event for %s not acknowledged", eventType)
	}

	var result *http.RetriesResult
	if !cloudevents.ResultAs(res, &result) {
		return fmt.Errorf("error decoding retries result %T: %+v", res, res)
	}

	var final *http.Result
	if !cloudevents.ResultAs(result.Result, &final) {
		return fmt.Errorf("error decoding final result %T: %+v", res, res)
	}

	retriesString := ""
	if result.Retries > 0 {
		retriesString = fmt.Sprintf(" (%d attempts)", result.Retries)
	}

	logrus.Infof("Sent %s with status: %d%s", eventType, final.StatusCode, retriesString)

	return nil
}

func (e *Events) NewUser(user User) error {
	err := e.sendEvent(e.sender, "ponglehub.auth.user.add", &user)

	if err != nil {
		return fmt.Errorf("failed to send add event for %s: %+v", user.Name, err)
	}

	return nil
}

func (e *Events) UpdateUser(user User) error {
	err := e.sendEvent(e.sender, "ponglehub.auth.user.update", &user)

	if err != nil {
		return fmt.Errorf("failed to send update event for %s: %+v", user.Name, err)
	}

	return nil
}

func (e *Events) DeleteUser(user User) error {
	err := e.sendEvent(e.sender, "ponglehub.auth.user.delete", &user)

	if err != nil {
		return fmt.Errorf("failed to send delete event for %s: %+v", user.Name, err)
	}

	return nil
}
