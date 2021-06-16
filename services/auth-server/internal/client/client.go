package client

import (
	"context"
)

// AuthClient - wrapper for database interactions

type AuthClient interface {
	AddUser(ctx context.Context, user User) (string, error)
	UpdateUser(ctx context.Context, id string, user User, verified *bool) (bool, error)
	GetUser(ctx context.Context, id string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	ListUsers(ctx context.Context) ([]*User, error)
	DeleteUser(ctx context.Context, id string) (bool, error)
}

// User - a struct representing a user
type User struct {
	ID       string
	Email    string
	Name     string
	Password string
	Verified bool
}
