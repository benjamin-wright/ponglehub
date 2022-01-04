module ponglehub.co.uk/events/broker

go 1.16

require (
	github.com/cloudevents/sdk-go/v2 v2.7.0
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	k8s.io/apimachinery v0.23.1
	k8s.io/client-go v0.23.1
	ponglehub.co.uk/events/recorder v1.0.0
	ponglehub.co.uk/lib/events v1.0.0
)

replace ponglehub.co.uk/lib/events => ./../../libraries/golang/events

replace ponglehub.co.uk/events/recorder => ./../event-recorder
