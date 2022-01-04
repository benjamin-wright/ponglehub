package nats

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	cenats "github.com/cloudevents/sdk-go/protocol/nats/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/cloudevents/sdk-go/v2/protocol/http"
	"github.com/sirupsen/logrus"
)

type Events struct {
	sender cloudevents.Client
	source string
}

type EventHandler func(ctx context.Context, event event.Event)

func New(brokerEnv string, subject string, source string) (*Events, error) {
	brokerUrl, ok := os.LookupEnv(brokerEnv)
	if !ok {
		logrus.Fatalf("Environment Variable %s not found", brokerEnv)
	}

	p, err := cenats.NewSender(brokerUrl, subject, cenats.NatsOptions())
	if err != nil {
		log.Fatalf("Failed to create nats protocol, %s", err.Error())
	}

	defer p.Close(context.Background())

	client, err := cloudevents.NewClient(p, cloudevents.WithTimeNow(), cloudevents.WithUUIDs())
	if err != nil {
		return nil, fmt.Errorf("error creating cloudevents instance: %+v", err)
	}

	return &Events{
		sender: client,
		source: source,
	}, nil
}

func Listen(brokerEnv string, subject string, handler EventHandler) error {
	brokerUrl, ok := os.LookupEnv(brokerEnv)
	if !ok {
		logrus.Fatalf("Environment Variable %s not found", brokerEnv)
	}

	logrus.Infof("Connecting with url: %s", brokerUrl)

	consumer, err := cenats.NewConsumer(brokerUrl, subject, cenats.NatsOptions())
	if err != nil {
		return fmt.Errorf("failed to create nats consumer: %s", err.Error())
	}
	defer consumer.Close(context.Background())

	c, err := cloudevents.NewClient(consumer)
	if err != nil {
		return fmt.Errorf("failed to create cloudevents client: %s", err.Error())
	}

	go func() {
		err := c.StartReceiver(context.Background(), handler)

		if err != nil {
			logrus.Fatalf("Error in event listener: %+v", err)
		} else {
			logrus.Infof("Stopped event listener")
		}
	}()

	return nil
}

func (e *Events) Send(eventType string, data interface{}) error {
	event := cloudevents.NewEvent()
	event.SetType(eventType)
	event.SetSource(e.source)
	err := event.SetData(cloudevents.ApplicationJSON, data)
	if err != nil {
		return fmt.Errorf("failed to serialize event data: %+v", err)
	}

	ctx := cloudevents.ContextWithRetriesConstantBackoff(context.TODO(), time.Second, 20)
	res := e.sender.Send(ctx, event)

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
