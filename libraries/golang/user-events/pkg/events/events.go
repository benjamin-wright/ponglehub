package events

import (
	"context"
	"fmt"

	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/lib/events/pkg/nats"
)

type Events struct {
	events *nats.Events
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

func New(brokerEnv string, source string) (*Events, error) {
	events, err := nats.New(brokerEnv, "com.ponglehub.auth", source)
	if err != nil {
		return nil, err
	}

	return &Events{
		events: events,
	}, nil
}

func Listen(brokerEnv string, handler UserEventHandler) error {
	return nats.Listen(brokerEnv, "com.ponglehub.auth", func(ctx context.Context, event event.Event) {
		user := User{}
		err := event.DataAs(&user)
		if err != nil {
			logrus.Errorf("Failed to parse event data: %+v", err)
			return
		}

		handler(UserEvent{
			Type: event.Type(),
			User: user,
		})
	})
}

func (e *Events) NewUser(user User) error {
	err := e.events.Send("ponglehub.auth.user.add", &user)

	if err != nil {
		return fmt.Errorf("failed to send add event for %s: %+v", user.Name, err)
	}

	return nil
}

func (e *Events) UpdateUser(user User) error {
	err := e.events.Send("ponglehub.auth.user.update", &user)

	if err != nil {
		return fmt.Errorf("failed to send update event for %s: %+v", user.Name, err)
	}

	return nil
}

func (e *Events) DeleteUser(user User) error {
	err := e.events.Send("ponglehub.auth.user.delete", &user)

	if err != nil {
		return fmt.Errorf("failed to send delete event for %s: %+v", user.Name, err)
	}

	return nil
}

func (e *Events) SetUser(user User) error {
	err := e.events.Send("ponglehub.auth.user.set", &user)

	if err != nil {
		return fmt.Errorf("failed to send set user event for %s: %+v", user.Name, err)
	}

	return nil
}

func (e *Events) SetUserAck(user User) error {
	err := e.events.Send("ponglehub.auth.user.set.ack", &user)

	if err != nil {
		return fmt.Errorf("failed to send set user ack event for %s: %+v", user.Name, err)
	}

	return nil
}
