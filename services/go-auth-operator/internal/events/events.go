package events

import (
	"context"
	"fmt"
	"time"

	cenats "github.com/cloudevents/sdk-go/protocol/nats/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/auth/auth-operator/internal/client"
)

type Events struct {
	p      *cenats.Sender
	sender cloudevents.Client
}

func New(natsUrl string, subject string) (*Events, error) {
	p, err := cenats.NewSender(natsUrl, subject, cenats.NatsOptions())
	if err != nil {
		return nil, fmt.Errorf("error connecting to NATS: %+v", err)
	}

	client, err := cloudevents.NewClient(p)
	if err != nil {
		return nil, fmt.Errorf("error creating cloudevents instance: %+v", err)
	}

	return &Events{
		p:      p,
		sender: client,
	}, nil
}

type NewUserData struct {
	MetaName string `json:"meta_name"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func sendEvent(sender cloudevents.Client, eventType string, data interface{}) error {
	event := cloudevents.NewEvent()
	event.SetID(uuid.New().String())
	event.SetType(eventType)
	event.SetTime(time.Now())
	event.SetSource("auth-operator")
	err := event.SetData("application/json", data)
	if err != nil {
		return fmt.Errorf("failed to serialize event data: %+v", err)
	}

	if res := sender.Send(context.TODO(), event); cloudevents.IsUndelivered(res) {
		return fmt.Errorf("failed to send event: %v", res.Error())
	}

	return nil
}

func (e *Events) NewUser(user *client.AuthUser) error {
	logrus.Infof("Sending new user event for %s", user.Name)

	err := sendEvent(e.sender, "ponglehub.auth.user.add", &NewUserData{
		MetaName: user.ObjectMeta.Name,
		Name:     user.Spec.Name,
		Email:    user.Spec.Email,
		Password: user.Spec.Password,
	})

	if err != nil {
		return fmt.Errorf("failed to send add event for %s: %+v", user.Name, err)
	}

	return nil
}

func (e *Events) UpdateUser(user *client.AuthUser) error {
	logrus.Infof("Sending update user event for %s", user.Name)

	err := sendEvent(e.sender, "ponglehub.auth.user.add", &NewUserData{
		MetaName: user.ObjectMeta.Name,
		Name:     user.Spec.Name,
		Email:    user.Spec.Email,
		Password: user.Spec.Password,
	})

	if err != nil {
		return fmt.Errorf("failed to send update event for %s: %+v", user.Name, err)
	}

	return nil
}

type DeleteUserData struct {
	MetaName string `json:"meta_name"`
	Name     string `json:"name"`
}

func (e *Events) DeleteUser(user *client.AuthUser) error {
	logrus.Infof("Sending delete user event for %s", user.Name)
	err := sendEvent(e.sender, "ponglehub.auth.user.add", &DeleteUserData{
		MetaName: user.ObjectMeta.Name,
		Name:     user.Spec.Name,
	})

	if err != nil {
		return fmt.Errorf("failed to send delete event for %s: %+v", user.Name, err)
	}

	return nil
}
