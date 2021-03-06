package client

import (
	"context"
	"fmt"
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
func (a *AuthClient) AddUser(ctx context.Context, user User) error {
	_, err := a.conn.Exec(
		ctx,
		"INSERT INTO user (id, name, email, password, verified) VALUES ($1, $2, $3, $4, False)",
		user.ID,
		user.Name,
		user.Email,
		user.Password,
	)

	return err
}

// GetUser - retrieve a user from the database
func (a *AuthClient) GetUser(ctx context.Context, name string) (*User, error) {
	rows, err := a.conn.Query(ctx, "SELECT id, name, email, password, verified FROM users WHERE name = $1", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Verified); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if len(users) == 0 {
		return nil, nil
	}

	if len(users) > 1 {
		return nil, fmt.Errorf("Failed to find user %s: Expected 1 user, received %d", name, len(users))
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

func (a *AuthClient) DeleteUser(ctx context.Context, id string) error {
	_, err := a.conn.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
}
