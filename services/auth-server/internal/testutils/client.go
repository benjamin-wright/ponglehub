package testutils

import (
	"context"
	"errors"

	pgx "github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"ponglehub.co.uk/auth/auth-server/internal/client"
	"ponglehub.co.uk/auth/auth-server/internal/migrations"
	"ponglehub.co.uk/lib/postgres/pkg/connect"
)

type TestClient struct {
	conn      *pgx.Conn
	targetCfg connect.ConnectConfig
	database  string
}

// New - Create a new AuthClient instance
func NewClient(database string) (*TestClient, error) {
	targetCfg, err := connect.ConfigFromEnv()
	if err != nil {
		logrus.Fatalf("Failed to load target config from environment: %+v", err)
	}

	adminCfg, err := connect.AdminFromEnv()
	if err != nil {
		logrus.Fatalf("Failed to load admin config from environment: %+v", err)
	}

	targetCfg.Database = database

	err = migrations.Migrate(targetCfg, adminCfg)
	if err != nil {
		return nil, err
	}

	conn, err := connect.Connect(targetCfg)
	if err != nil {
		return nil, err
	}

	return &TestClient{conn: conn, targetCfg: targetCfg, database: database}, nil
}

func (c *TestClient) TargetConfig() connect.ConnectConfig {
	return c.targetCfg
}

func (c *TestClient) Close() error {
	return c.conn.Close(context.Background())
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
