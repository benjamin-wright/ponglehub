module ponglehub.co.uk/lib/user-events

go 1.16

require (
	github.com/cloudevents/sdk-go/v2 v2.7.0
	github.com/sirupsen/logrus v1.8.1
	ponglehub.co.uk/lib/events v1.0.0
)

replace ponglehub.co.uk/lib/events => ./../events
