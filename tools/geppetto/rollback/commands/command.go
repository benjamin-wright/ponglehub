package commands

// Command an interface that represents a task in the rolling back of a repo
type Command interface {
	Run() error
	Name() string
}
