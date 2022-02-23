module ponglehub.co.uk/events/responder

go 1.17

require (
	github.com/go-redis/redis/v8 v8.11.4
	github.com/stretchr/testify v1.5.1
	ponglehub.co.uk/lib/events v0.0.0-00010101000000-000000000000
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/google/uuid v1.1.1 // indirect
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v0.0.0-20180701023420-4b7aa43c6742 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.uber.org/atomic v1.4.0 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.10.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/cloudevents/sdk-go/v2 v2.8.0
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/sirupsen/logrus v1.8.1
	golang.org/x/sys v0.0.0-20210615035016-665e8c7367d1 // indirect
)

replace ponglehub.co.uk/lib/events => ./../../libraries/golang/events
