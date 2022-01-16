module ponglehub.co.uk/games/naughts-and-crosses

go 1.17

require github.com/sirupsen/logrus v1.8.1

require (
	github.com/cloudevents/sdk-go/v2 v2.7.0
	ponglehub.co.uk/lib/events v1.0.0
	ponglehub.co.uk/lib/postgres v1.0.0
)

require (
	github.com/google/uuid v1.1.1 // indirect
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v0.0.0-20180701023420-4b7aa43c6742 // indirect
	go.uber.org/atomic v1.4.0 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.10.0 // indirect
	golang.org/x/sys v0.0.0-20210615035016-665e8c7367d1 // indirect
)

replace ponglehub.co.uk/lib/events => ./../../libraries/golang/events
replace ponglehub.co.uk/lib/postgres => ./../../libraries/golang/postgres
