package commands

// Command an interface that represents a task in the building of a repo
type Command interface {
	Run() (bool, error)
	Stage() string
}

type GenericCommand struct {
	runner func() (bool, error)
	stage  string
}

func CreateGeneric(stage string, runner func() (bool, error)) GenericCommand {
	return GenericCommand{
		runner: runner,
		stage:  stage,
	}
}

func (cmd GenericCommand) Run() (bool, error) {
	return cmd.runner()
}

func (cmd GenericCommand) Stage() string {
	return cmd.stage
}
