package client

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

// AuthClient - wrapper for database interactions
type AuthClient struct {
	conn *pgx.Conn
}

// AuthClientConfig - creds and config for creating a database connection
type AuthClientConfig struct {
	Username string
	Password string
	Host     string
	Port     int16
}

// New - Create a new AuthClient instance
func New(ctx context.Context, config *AuthClientConfig) (*AuthClient, error) {
	pgxConfig, err := pgx.ParseConfig(fmt.Sprintf("postgres://%s:%s@%s:%d", config.Username, config.Password, config.Host, config.Port))
	if err != nil {
		return nil, err
	}

	conn, err := pgx.ConnectConfig(ctx, pgxConfig)
	if err != nil {
		return nil, err
	}

	return &AuthClient{
		conn: conn,
	}, nil
}

func (a *AuthClient) Close(ctx context.Context) error {
	// Close - Remember to call this when you're done with the client
	return a.Close(ctx)
}
