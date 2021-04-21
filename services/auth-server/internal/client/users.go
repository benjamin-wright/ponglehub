package client

import (
	"context"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/sirupsen/logrus"
)

// User - a struct representing a user
type User struct {
	ID       string
	Email    string
	Name     string
	Password string
	Verified bool
}

// AddUser - Add a user to the database
func (a *AuthClient) AddUser(ctx context.Context, user User) (bool, error) {
	_, err := a.conn.Exec(
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
				return false, nil
			}
		}

		return false, err
	}

	return true, nil
}

// GetUser - retrieve a user from the database
func (a *AuthClient) GetUser(ctx context.Context, id string) (*User, error) {
	rows, err := a.conn.Query(ctx, "SELECT name, email, password, verified FROM users WHERE id = $1", id)
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

// ListUsers - Returns a list of all the users in the system
func (a *AuthClient) ListUsers(ctx context.Context) ([]*User, error) {
	rows, err := a.conn.Query(ctx, "SELECT id, name, email, password, verified FROM users")
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
func (a *AuthClient) DeleteUser(ctx context.Context, id string) (bool, error) {
	cmd, err := a.conn.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
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
