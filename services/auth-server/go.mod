module ponglehub.co.uk/auth/auth-server

go 1.16

require (
	github.com/gin-gonic/gin v1.7.1
	github.com/google/uuid v1.2.0
	github.com/jackc/pgconn v1.8.1
	github.com/jackc/pgerrcode v0.0.0-20201024163028-a0d42d470451
	github.com/jackc/pgx/v4 v4.11.0
	github.com/sirupsen/logrus v1.8.1
	golang.org/x/crypto v0.0.0-20210421170649-83a5a9bb288b
	gopkg.in/yaml.v2 v2.4.0

	ponglehub.co.uk/auth/db-init v1.0.0
)

replace ponglehub.co.uk/auth/db-init => ./../db-init
