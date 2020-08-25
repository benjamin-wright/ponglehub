package commands

// Command an interface that represents a task in the building of a repo
type Command interface {
	Run() (bool, error)
	Stage() string
}
