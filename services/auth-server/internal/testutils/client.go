package testutils

import (
	"context"
	"errors"
	"fmt"
	"os"

	pgx "github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"ponglehub.co.uk/auth/auth-server/internal/client"
	"ponglehub.co.uk/auth/auth-server/internal/migrations"
)

type TestClient struct {
	conn     *pgx.Conn
	admin    *pgx.Conn
	database string
}

// New - Create a new AuthClient instance
func NewClient(database string) (*TestClient, error) {
	host, ok := os.LookupEnv("DB_HOST")
	if !ok {
		logrus.Fatal("Environment Variable DB_HOST not found")
	}

	username, ok := os.LookupEnv("DB_USER")
	if !ok {
		username = "authserver"
	}

	password, ok := os.LookupEnv("DB_PASS")
	if !ok {
		logrus.Fatal("Enrivonment Variable DB_PASS not found")
	}

	certsDir, ok := os.LookupEnv("DB_CERTS")
	if !ok {
		logrus.Fatal("Enrivonment Variable DB_CERTS not found")
	}

	err := migrations.Migrate(host, username, password, database, certsDir)
	if err != nil {
		return nil, err
	}

	pgxConfig, err := pgx.ParseConfig(fmt.Sprintf("postgres://root@%s:26257/%s", host, database))
	if err != nil {
		return nil, err
	}

	conn, err := pgx.ConnectConfig(context.TODO(), pgxConfig)
	if err != nil {
		return nil, err
	}

	return &TestClient{conn: conn, database: database}, nil
}

func (c *TestClient) Drop() error {
	c.conn.Close(context.Background())

	host, ok := os.LookupEnv("DB_HOST")
	if !ok {
		return errors.New("DB_HOST env var must be provided")
	}

	config, err := pgx.ParseConfig(fmt.Sprintf("postgres://root@%s:26257", host))
	if err != nil {
		return err
	}

	conn, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		return fmt.Errorf("error connecting to the database: %+v", err)
	}
	defer conn.Close(context.Background())

	_, err = conn.Exec(
		context.TODO(),
		fmt.Sprintf("DROP DATABASE %s;", c.database),
	)

	return err
}

func (a *TestClient) Reset() error {
	_, err := a.conn.Exec(context.Background(), "DELETE FROM users")
	return err
}

func (a *TestClient) AddUser(id string, name string, email string, password string, verified bool) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = a.conn.Exec(
		context.Background(),
		"INSERT INTO users (id, name, email, password, verified) VALUES ($1, $2, $3, $4, $5)",
		id,
		name,
		email,
		hashed,
		verified,
	)

	return err
}

func (a *TestClient) GetUser(id string) (*client.User, error) {
	var user client.User

	rows, err := a.conn.Query(
		context.Background(),
		"SELECT name, email, password, verified FROM users WHERE id = $1",
		id,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if hasResult := rows.Next(); !hasResult {
		return nil, errors.New("failed to fetch user, returned less than one row")
	}

	rows.Scan(&user.Name, &user.Email, &user.Password, &user.Verified)

	return &user, nil
}

func (a *TestClient) ListUserIds() ([]string, error) {
	rows, err := a.conn.Query(
		context.Background(),
		"SELECT id FROM users",
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := []string{}

	for rows.Next() {
		var id string
		err = rows.Scan(&id)

		if err != nil {
			return nil, err
		}

		ids = append(ids, id)
	}

	return ids, nil
}
