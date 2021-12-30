package http

import (
	"context"
	"fmt"
	"os"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/cloudevents/sdk-go/v2/protocol/http"
	"github.com/sirupsen/logrus"
)

type Events struct {
	ctx    context.Context
	sender cloudevents.Client
	source string
}

type EventHandler func(ctx context.Context, event event.Event)

func New(brokerEnv string, source string) (*Events, error) {
	brokerUrl, ok := os.LookupEnv(brokerEnv)
	if !ok {
		logrus.Fatalf("Environment Variable %s not found", brokerEnv)
	}

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
		source: source,
	}, nil
}

func Listen(handler EventHandler) (context.CancelFunc, error) {
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
		err := client.StartReceiver(ctx, handler)

		if err != nil && ctx.Err() == nil {
			logrus.Fatalf("Error in event listener: %+v", err)
		} else {
			logrus.Infof("Stopped event listener")
		}
	}()

	return cancel, nil
}

func (e *Events) Send(eventType string, data interface{}) error {
	event := cloudevents.NewEvent()
	event.SetType(eventType)
	event.SetSource(e.source)
	err := event.SetData(cloudevents.ApplicationJSON, data)
	if err != nil {
		return fmt.Errorf("failed to serialize event data: %+v", err)
	}

	ctx := cloudevents.ContextWithRetriesConstantBackoff(e.ctx, time.Second, 20)
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
