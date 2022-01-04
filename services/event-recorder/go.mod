module ponglehub.co.uk/events/recorder

go 1.17

require (
	github.com/cloudevents/sdk-go/v2 v2.7.0
	github.com/gin-gonic/gin v1.7.7
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.5.1
	ponglehub.co.uk/lib/events v1.0.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.13.0 // indirect
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/go-playground/validator/v10 v10.4.1 // indirect
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/google/uuid v1.1.1 // indirect
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v0.0.0-20180701023420-4b7aa43c6742 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/ugorji/go/codec v1.1.7 // indirect
	go.uber.org/atomic v1.4.0 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.10.0 // indirect
	golang.org/x/crypto v0.0.0-20210616213533-5ff15b29337e // indirect
	golang.org/x/sys v0.0.0-20210615035016-665e8c7367d1 // indirect
	google.golang.org/protobuf v1.23.0 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
)

replace ponglehub.co.uk/lib/events => ./../../libraries/golang/events
