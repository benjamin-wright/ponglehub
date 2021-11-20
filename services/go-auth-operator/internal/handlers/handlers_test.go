package handlers

import (
	"errors"
	"io"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"ponglehub.co.uk/auth/auth-operator/internal/events"
)

type MockEvents struct {
	mock.Mock
}

func (m *MockEvents) NewUser(user events.User) error {
	return m.Called(user).Error(0)
}

func (m *MockEvents) UpdateUser(user events.User) error {
	return m.Called(user).Error(0)
}

func (m *MockEvents) DeleteUser(user events.User) error {
	return m.Called(user).Error(0)
}

type MockUsers struct {
	mock.Mock
}

func (m *MockUsers) Status(user events.User) (events.User, error) {
	args := m.Called(user)
	return args.Get(0).(events.User), args.Error(1)
}

func (m *MockUsers) Update(user events.User) (events.User, error) {
	args := m.Called(user)
	return args.Get(0).(events.User), args.Error(1)
}

func (m *MockUsers) Get(name string) (events.User, error) {
	args := m.Called(name)
	return args.Get(0).(events.User), args.Error(1)
}

func makeUser(name string, email string, password string, id string, pending bool) events.User {
	return events.User{
		Name:     name,
		Username: name,
		Email:    email,
		Password: password,
		ID:       id,
		Pending:  pending,
	}
}

func TestAddHander(t *testing.T) {
	logrus.SetOutput(io.Discard)

	for _, test := range []struct {
		Name  string
		Prep  func(*MockEvents, *MockUsers)
		Input events.User
	}{
		{
			Name:  "New user",
			Input: makeUser("name", "email", "pass", "", false),
			Prep: func(mockEvents *MockEvents, users *MockUsers) {
				users.On("Status", makeUser("name", "email", "pass", "", true)).Return(events.User{}, nil)
				mockEvents.On("NewUser", makeUser("name", "email", "pass", "", true)).Return(nil)
			},
		},
		{
			Name:  "Failed to create user event",
			Input: makeUser("name", "email", "pass", "", false),
			Prep: func(mockEvents *MockEvents, users *MockUsers) {
				users.On("Status", makeUser("name", "email", "pass", "", true)).Return(events.User{}, nil)
				mockEvents.On("NewUser", makeUser("name", "email", "pass", "", true)).Return(errors.New("boom"))
			},
		},
		{
			Name:  "Failed setting user status",
			Input: makeUser("name", "email", "pass", "", false),
			Prep: func(mockEvents *MockEvents, users *MockUsers) {
				users.On("Status", makeUser("name", "email", "pass", "", true)).Return(events.User{}, errors.New("boom"))
			},
		},
		{
			Name:  "Already pending",
			Input: makeUser("name", "email", "pass", "", true),
			Prep:  func(mockEvents *MockEvents, users *MockUsers) {},
		},
		{
			Name:  "Got id",
			Input: makeUser("name", "email", "pass", "1234", false),
			Prep:  func(mockEvents *MockEvents, users *MockUsers) {},
		},
	} {
		t.Run(test.Name, func(u *testing.T) {
			mockEvents := new(MockEvents)
			mockUsers := new(MockUsers)
			test.Prep(mockEvents, mockUsers)

			handlers := Handler{events: mockEvents, users: mockUsers}

			handlers.AddUser(test.Input)

			mockEvents.AssertExpectations(u)
			mockUsers.AssertExpectations(u)
		})
	}
}

func TestUpdateHander(t *testing.T) {
	logrus.SetOutput(io.Discard)

	for _, test := range []struct {
		Name    string
		Prep    func(*MockEvents, *MockUsers)
		OldUser events.User
		NewUser events.User
	}{
		{
			Name:    "Updated user",
			OldUser: makeUser("name", "email", "password", "", false),
			NewUser: makeUser("name", "new-email", "new-password", "", false),
			Prep: func(mockEvents *MockEvents, users *MockUsers) {
				users.On("Status", makeUser("name", "new-email", "new-password", "", true)).Return(events.User{}, nil)
				mockEvents.On("UpdateUser", makeUser("name", "new-email", "new-password", "", true)).Return(nil)
			},
		},
		{
			Name:    "Failed to send update event",
			OldUser: makeUser("name", "email", "password", "", false),
			NewUser: makeUser("name", "new-email", "new-password", "", false),
			Prep: func(mockEvents *MockEvents, users *MockUsers) {
				users.On("Status", makeUser("name", "new-email", "new-password", "", true)).Return(events.User{}, nil)
				mockEvents.On("UpdateUser", makeUser("name", "new-email", "new-password", "", true)).Return(errors.New("boom"))
			},
		},
		{
			Name:    "fail setting status",
			OldUser: makeUser("name", "email", "password", "", false),
			NewUser: makeUser("name", "new-email", "new-password", "", false),
			Prep: func(mockEvents *MockEvents, users *MockUsers) {
				users.On("Status", makeUser("name", "new-email", "new-password", "", true)).Return(events.User{}, errors.New("oops"))
			},
		},
		{
			Name:    "No spec change",
			OldUser: makeUser("name", "email", "password", "", false),
			NewUser: makeUser("name", "email", "password", "", false),
			Prep:    func(mockEvents *MockEvents, users *MockUsers) {},
		},
		{
			Name:    "Already changing",
			OldUser: makeUser("name", "email", "password", "", false),
			NewUser: makeUser("name", "new-email", "new-password", "", true),
			Prep:    func(mockEvents *MockEvents, users *MockUsers) {},
		},
	} {
		t.Run(test.Name, func(u *testing.T) {
			mockEvents := new(MockEvents)
			mockUsers := new(MockUsers)
			test.Prep(mockEvents, mockUsers)

			handlers := Handler{events: mockEvents, users: mockUsers}

			handlers.UpdateUser(test.OldUser, test.NewUser)

			mockEvents.AssertExpectations(u)
			mockUsers.AssertExpectations(u)
		})
	}
}

func TestDeleteHander(t *testing.T) {
	logrus.SetOutput(io.Discard)

	for _, test := range []struct {
		Name  string
		Prep  func(*MockEvents, *MockUsers)
		Input events.User
	}{
		{
			Name:  "Delete user",
			Input: makeUser("name", "email", "pass", "", false),
			Prep: func(events *MockEvents, users *MockUsers) {
				events.On("DeleteUser", makeUser("name", "email", "pass", "", false)).Return(nil)
			},
		},
		{
			Name:  "Error deleting user",
			Input: makeUser("name", "email", "pass", "", false),
			Prep: func(events *MockEvents, users *MockUsers) {
				events.On("DeleteUser", makeUser("name", "email", "pass", "", false)).Return(errors.New("boom"))
			},
		},
	} {
		t.Run(test.Name, func(u *testing.T) {
			mockEvents := new(MockEvents)
			mockUsers := new(MockUsers)
			test.Prep(mockEvents, mockUsers)

			handlers := Handler{events: mockEvents, users: mockUsers}

			handlers.DeleteUser(test.Input)

			mockEvents.AssertExpectations(u)
			mockUsers.AssertExpectations(u)
		})
	}
}
