package client

import (
	"context"

	"github.com/google/uuid"
)

type MockClient struct {
	users []User
}

func NewMockClient(users []User) *MockClient {
	return &MockClient{users}
}

// AddUser - Add a user to the database
func (p *MockClient) AddUser(ctx context.Context, user User) (string, error) {
	for _, u := range p.users {
		if u.ID == user.ID || u.Name == user.Name || u.Email == user.Email {
			return "", nil
		}
	}

	user.ID = uuid.New().String()
	p.users = append(p.users, user)

	return user.ID, nil
}

// UpdateUser - Update an existing user
func (p *MockClient) UpdateUser(ctx context.Context, id string, user User, verified *bool) (bool, error) {
	for _, u := range p.users {
		if u.ID == id {
			if user.Email != "" {
				u.Email = user.Email
			}

			if user.Name != "" {
				u.Name = user.Name
			}

			if user.Password != "" {
				u.Password = user.Password
			}

			if verified != nil {
				u.Verified = *verified
			}

			return true, nil
		}
	}

	return false, nil
}

// GetUser - retrieve a user from the database
func (p *MockClient) GetUser(ctx context.Context, id string) (*User, error) {
	for _, u := range p.users {
		if u.ID == id {
			return &u, nil
		}
	}

	return nil, nil
}

// GetUser - retrieve a user from the database
func (p *MockClient) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	for _, u := range p.users {
		if u.Email == email {
			return &u, nil
		}
	}

	return nil, nil
}

// ListUsers - Returns a list of all the users in the system
func (p *MockClient) ListUsers(ctx context.Context) ([]*User, error) {
	result := []*User{}

	for _, u := range p.users {
		result = append(result, &u)
	}

	return result, nil
}

// DeleteUser - delete a user
func (p *MockClient) DeleteUser(ctx context.Context, id string) (bool, error) {
	for idx, u := range p.users {
		if u.ID == id {
			l := len(p.users)
			p.users[idx] = p.users[l-1]
			p.users = p.users[:l-1]
			return true, nil
		}
	}

	return false, nil
}
