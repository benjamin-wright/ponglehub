package client

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	pgx "github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
)

type PostgresClient struct {
	conn *pgx.Conn
}

// AuthClientConfig - creds and config for creating a database connection
type PostgresClientConfig struct {
	Username string
	Password string
	Host     string
	Port     int16
	Database string
}

// New - Create a new AuthClient instance
func NewPostgresClient(ctx context.Context, config *PostgresClientConfig) (*PostgresClient, error) {
	pgxConfig, err := pgx.ParseConfig(fmt.Sprintf("postgres://%s:%s@%s:%d/%s", config.Username, config.Password, config.Host, config.Port, config.Database))
	if err != nil {
		return nil, err
	}

	conn, err := pgx.ConnectConfig(ctx, pgxConfig)
	if err != nil {
		return nil, err
	}

	return &PostgresClient{
		conn: conn,
	}, nil
}

// Close - Remember to call this when you're done with the client
func (p *PostgresClient) Close(ctx context.Context) error {
	return p.conn.Close(ctx)
}

// AddUser - Add a user to the database
func (p *PostgresClient) AddUser(ctx context.Context, user User) (string, error) {
	_, err := p.conn.Exec(
		ctx,
		"INSERT INTO users (name, email, password, verified) VALUES ($1, $2, $3, False)",
		user.Name,
		user.Email,
		user.Password,
	)

	if err != nil {
		if err, ok := err.(*pgconn.PgError); ok {
			if err.Code == pgerrcode.UniqueViolation {
				logrus.Warnf("Failed to create duplicate user: %s[%s]", user.Name, user.Email)
				return "", nil
			}
		}

		return "", err
	}

	rows, err := p.conn.Query(ctx, "SELECT id from users WHERE name = $1", user.Name)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	id := ""
	if hasNext := rows.Next(); !hasNext {
		return "", errors.New("Failed to create user: no id response")
	}

	rows.Scan(&id)

	return id, nil
}

// UpdateUser - Update an existing user
func (p *PostgresClient) UpdateUser(ctx context.Context, id string, user User, verified *bool) (bool, error) {
	queryParts := []string{}
	queryArgs := []interface{}{}

	term := 1

	if user.Email != "" {
		queryParts = append(queryParts, fmt.Sprintf("email = $%d", term))
		queryArgs = append(queryArgs, user.Email)

		term += 1
	}

	if user.Name != "" {
		queryParts = append(queryParts, fmt.Sprintf("name = $%d", term))
		queryArgs = append(queryArgs, user.Name)

		term += 1
	}

	if user.Password != "" {
		queryParts = append(queryParts, fmt.Sprintf("password = $%d", term))
		queryArgs = append(queryArgs, user.Password)

		term += 1
	}

	if verified != nil {
		queryParts = append(queryParts, fmt.Sprintf("verified = $%d", term))
		queryArgs = append(queryArgs, &verified)

		term += 1
	}

	query := fmt.Sprintf("UPDATE users SET %s WHERE id = $%d", strings.Join(queryParts, ", "), term)
	queryArgs = append(queryArgs, id)

	res, err := p.conn.Exec(
		ctx,
		query,
		queryArgs...,
	)
	if err != nil {
		return false, err
	}

	if res.RowsAffected() != 1 {
		return false, nil
	}

	return true, nil
}

// GetUser - retrieve a user from the database
func (p *PostgresClient) GetUser(ctx context.Context, id string) (*User, error) {
	rows, err := p.conn.Query(ctx, "SELECT name, email, password, verified FROM users WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Name, &user.Email, &user.Password, &user.Verified); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if len(users) == 0 {
		return nil, nil
	}

	if len(users) > 1 {
		return nil, fmt.Errorf("Failed to find user %s: Expected 1 user, received %d", id, len(users))
	}

	return &users[0], nil
}

// GetUser - retrieve a user from the database
func (p *PostgresClient) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	rows, err := p.conn.Query(ctx, "SELECT id, name, password, verified FROM users WHERE email = $1", email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Password, &user.Verified); err != nil {
			return nil, err
		}

		user.Email = email

		users = append(users, user)
	}

	if len(users) == 0 {
		return nil, nil
	}

	if len(users) > 1 {
		return nil, fmt.Errorf("Failed to find user from email \"%s\": Expected 1 user, received %d", email, len(users))
	}

	return &users[0], nil
}

// ListUsers - Returns a list of all the users in the system
func (p *PostgresClient) ListUsers(ctx context.Context) ([]*User, error) {
	rows, err := p.conn.Query(ctx, "SELECT id, name, email, password, verified FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*User{}

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Verified); err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

// DeleteUser - delete a user
func (p *PostgresClient) DeleteUser(ctx context.Context, id string) (bool, error) {
	cmd, err := p.conn.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return false, err
	}

	numRows := cmd.RowsAffected()
	if numRows == 0 {
		return false, nil
	}

	logrus.Infof("Deleted %d users", numRows)

	return true, nil
}
