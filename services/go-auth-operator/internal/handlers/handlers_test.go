package handlers

import (
	"errors"
	"io"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"ponglehub.co.uk/auth/auth-operator/internal/client"
)

type MockEvents struct {
	mock.Mock
}

func (m *MockEvents) NewUser(user *client.AuthUser) error {
	return m.Called(user).Error(0)
}

func (m *MockEvents) UpdateUser(user *client.AuthUser) error {
	return m.Called(user).Error(0)
}

func (m *MockEvents) DeleteUser(user *client.AuthUser) error {
	return m.Called(user).Error(0)
}

type MockUsers struct {
	mock.Mock
}

func (m *MockUsers) SetStatus(user *client.AuthUser, opts v1.UpdateOptions) error {
	return m.Called(user, opts).Error(0)
}

func (m *MockUsers) Get(name string, opts v1.GetOptions) (*client.AuthUser, error) {
	args := m.Called(name, opts)
	return args.Get(0).(*client.AuthUser), args.Error(1)
}

func makeUser(name string, email string, password string, id string, pending bool) *client.AuthUser {
	return &client.AuthUser{
		ObjectMeta: v1.ObjectMeta{
			Name: name,
		},
		Spec: client.AuthUserSpec{
			Name:     name,
			Email:    email,
			Password: password,
		},
		Status: client.AuthUserStatus{
			ID:      id,
			Pending: pending,
		},
	}
}

func TestAddHander(t *testing.T) {
	logrus.SetOutput(io.Discard)

	for _, test := range []struct {
		Name  string
		Prep  func(*MockEvents, *MockUsers)
		Input *client.AuthUser
	}{
		{
			Name: "No input user",
			Prep: func(events *MockEvents, users *MockUsers) {},
		},
		{
			Name:  "New user",
			Input: makeUser("name", "email", "pass", "", false),
			Prep: func(events *MockEvents, users *MockUsers) {
				users.On("SetStatus", makeUser("name", "email", "pass", "", true), v1.UpdateOptions{}).Return(nil)
				events.On("NewUser", makeUser("name", "email", "pass", "", true)).Return(nil)
			},
		},
		{
			Name:  "Failed to create user event",
			Input: makeUser("name", "email", "pass", "", false),
			Prep: func(events *MockEvents, users *MockUsers) {
				users.On("SetStatus", makeUser("name", "email", "pass", "", true), v1.UpdateOptions{}).Return(nil)
				events.On("NewUser", makeUser("name", "email", "pass", "", true)).Return(errors.New("boom"))
				users.On("Get", "name", v1.GetOptions{}).Return(makeUser("name", "email", "pass", "", true), nil)
				users.On("SetStatus", makeUser("name", "email", "pass", "", false), v1.UpdateOptions{}).Return(nil)
			},
		},
		{
			Name:  "Failed setting user status",
			Input: makeUser("name", "email", "pass", "", false),
			Prep: func(events *MockEvents, users *MockUsers) {
				users.On("SetStatus", makeUser("name", "email", "pass", "", true), v1.UpdateOptions{}).Return(errors.New("boom"))
			},
		},
		{
			Name:  "Already pending",
			Input: makeUser("name", "email", "pass", "", true),
			Prep:  func(events *MockEvents, users *MockUsers) {},
		},
		{
			Name:  "Got id",
			Input: makeUser("name", "email", "pass", "1234", false),
			Prep:  func(events *MockEvents, users *MockUsers) {},
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
		OldUser *client.AuthUser
		NewUser *client.AuthUser
	}{
		{
			Name: "No input user",
			Prep: func(events *MockEvents, users *MockUsers) {},
		},
		{
			Name:    "Updated user",
			OldUser: makeUser("name", "email", "password", "", false),
			NewUser: makeUser("name", "new-email", "new-password", "", false),
			Prep: func(events *MockEvents, users *MockUsers) {
				users.On("SetStatus", makeUser("name", "new-email", "new-password", "", true), v1.UpdateOptions{}).Return(nil)
				events.On("UpdateUser", makeUser("name", "new-email", "new-password", "", true)).Return(nil)
			},
		},
		{
			Name:    "Failed to send update event",
			OldUser: makeUser("name", "email", "password", "", false),
			NewUser: makeUser("name", "new-email", "new-password", "", false),
			Prep: func(events *MockEvents, users *MockUsers) {
				users.On("SetStatus", makeUser("name", "new-email", "new-password", "", true), v1.UpdateOptions{}).Return(nil)
				events.On("UpdateUser", makeUser("name", "new-email", "new-password", "", true)).Return(errors.New("boom"))
				users.On("Get", "name", v1.GetOptions{}).Return(makeUser("name", "email", "pass", "", true), nil)
				users.On("SetStatus", makeUser("name", "email", "pass", "", false), v1.UpdateOptions{}).Return(nil)
			},
		},
		{
			Name:    "fail setting status",
			OldUser: makeUser("name", "email", "password", "", false),
			NewUser: makeUser("name", "new-email", "new-password", "", false),
			Prep: func(events *MockEvents, users *MockUsers) {
				users.On("SetStatus", makeUser("name", "new-email", "new-password", "", true), v1.UpdateOptions{}).Return(errors.New("oops"))
			},
		},
		{
			Name:    "No spec change",
			OldUser: makeUser("name", "email", "password", "", false),
			NewUser: makeUser("name", "email", "password", "", false),
			Prep:    func(events *MockEvents, users *MockUsers) {},
		},
		{
			Name:    "Already changing",
			OldUser: makeUser("name", "email", "password", "", false),
			NewUser: makeUser("name", "new-email", "new-password", "", true),
			Prep:    func(events *MockEvents, users *MockUsers) {},
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
		Input *client.AuthUser
	}{
		{
			Name: "No input user",
			Prep: func(events *MockEvents, users *MockUsers) {},
		},
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
