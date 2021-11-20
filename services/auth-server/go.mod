module ponglehub.co.uk/auth/auth-server

go 1.16

require (
	github.com/gin-gonic/gin v1.7.1
	github.com/google/uuid v1.2.0
	github.com/jackc/pgconn v1.10.0
	github.com/jackc/pgerrcode v0.0.0-20201024163028-a0d42d470451
	github.com/jackc/pgx/v4 v4.13.0
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97
	gopkg.in/yaml.v2 v2.4.0 // indirect
	ponglehub.co.uk/lib/postgres v1.0.0
	ponglehub.co.uk/lib/user-events v1.0.0
)

replace ponglehub.co.uk/lib/postgres => ./../../libraries/golang/postgres

replace ponglehub.co.uk/lib/user-events => ./../../libraries/golang/user-events
