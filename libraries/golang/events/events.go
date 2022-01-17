package events

import (
	"context"
	"fmt"
	nethttp "net/http"
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

type EventsArgs struct {
	BrokerEnv string
	BrokerURL string
	Source    string
	Cookies   nethttp.CookieJar
}

type Error string

func (e Error) Error() string { return string(e) }

const UnauthorizedError = Error("unauthorized")

func New(args EventsArgs) (*Events, error) {
	brokerUrl := args.BrokerURL

	if args.BrokerEnv != "" {
		url, ok := os.LookupEnv(args.BrokerEnv)
		if !ok {
			return nil, fmt.Errorf("environment variable %s not found", args.BrokerEnv)
		}

		brokerUrl = url
	}

	if brokerUrl == "" {
		return nil, fmt.Errorf("no broker url found, provide either BrokerEnv or BrokerURL")
	}

	ctx := cloudevents.ContextWithTarget(context.Background(), brokerUrl)

	p, err := cloudevents.NewHTTP()
	if err != nil {
		return nil, fmt.Errorf("error connecting to NATS: %+v", err)
	}

	p.Client.Timeout = time.Second
	p.Client.Jar = args.Cookies

	client, err := cloudevents.NewClient(p, cloudevents.WithTimeNow(), cloudevents.WithUUIDs())
	if err != nil {
		return nil, fmt.Errorf("error creating cloudevents instance: %+v", err)
	}

	return &Events{
		ctx:    ctx,
		sender: client,
		source: args.Source,
	}, nil
}

func (e *Events) Proxy(event event.Event) error {
	ctx := cloudevents.ContextWithRetriesConstantBackoff(e.ctx, time.Second, 20)
	res := e.sender.Send(ctx, event)

	if cloudevents.IsUndelivered(res) {
		return fmt.Errorf("failed to send event: %v", res.Error())
	}

	var result *http.RetriesResult
	if !cloudevents.ResultAs(res, &result) {
		return fmt.Errorf("error decoding retries result %T: %+v", res, res)
	}

	var final *http.Result
	if !cloudevents.ResultAs(result.Result, &final) {
		return fmt.Errorf("error decoding final result %T: %+v", res, res)
	}

	if cloudevents.IsNACK(res) {
		if final.StatusCode == 401 {
			return UnauthorizedError
		} else {
			return fmt.Errorf("event for %s not acknowledged: %d", event.Type(), final.StatusCode)
		}
	}

	retriesString := ""
	if result.Retries > 0 {
		retriesString = fmt.Sprintf(" (%d attempts)", result.Retries)
	}

	logrus.Infof("Sent %s with status: %d%s", event.Type(), final.StatusCode, retriesString)

	return nil
}

func (e *Events) Send(eventType string, data interface{}, extensions ...map[string]interface{}) error {
	event := cloudevents.NewEvent()
	event.SetType(eventType)
	event.SetSource(e.source)

	for _, ext := range extensions {
		for key, value := range ext {
			event.SetExtension(key, value)
		}
	}

	err := event.SetData(cloudevents.ApplicationJSON, data)
	if err != nil {
		return fmt.Errorf("failed to serialize event data: %+v", err)
	}

	return e.Proxy(event)
}

type EventHandler func(ctx context.Context, event event.Event)

func Listen(port int, handler EventHandler) (context.CancelFunc, error) {
	p, err := cloudevents.NewHTTP(cloudevents.WithPort(port))
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

	logrus.Infof("Listening on port %d...", port)

	return cancel, nil
}
