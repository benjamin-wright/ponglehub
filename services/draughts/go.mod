module ponglehub.co.uk/games/draughts

go 1.17

replace ponglehub.co.uk/lib/events => ./../../libraries/golang/events

replace ponglehub.co.uk/lib/postgres => ./../../libraries/golang/postgres

replace ponglehub.co.uk/events/recorder => ./../event-recorder

require (
	github.com/cloudevents/sdk-go/v2 v2.7.0
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	ponglehub.co.uk/events/recorder v0.0.0-00010101000000-000000000000
	ponglehub.co.uk/lib/postgres v0.0.0-00010101000000-000000000000
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v0.0.0-20180701023420-4b7aa43c6742 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.uber.org/atomic v1.6.0 // indirect
	go.uber.org/multierr v1.5.0 // indirect
	go.uber.org/zap v1.13.0 // indirect
	golang.org/x/lint v0.0.0-20190930215403-16217165b5de // indirect
	golang.org/x/tools v0.0.0-20200103221440-774c71fcf114 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c // indirect
)

require (
	github.com/google/uuid v1.3.0
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.10.0 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.1.1 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/pgtype v1.8.1 // indirect
	github.com/jackc/pgx/v4 v4.13.0
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97 // indirect
	golang.org/x/sys v0.0.0-20210615035016-665e8c7367d1 // indirect
	golang.org/x/text v0.3.6 // indirect
	ponglehub.co.uk/lib/events v1.0.0
)
