package state

type State struct {
	Users []User
}

type User struct {
	ID       string
	Name     string
	Email    string
	Password string
	Verified bool
}

func New() *State {
	return &State{
		Users: []User{},
	}
}
