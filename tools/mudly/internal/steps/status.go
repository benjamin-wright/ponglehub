package steps

type CommandResult int

const (
	COMMAND_SUCCESS CommandResult = iota
	COMMAND_ERROR
	COMMAND_SKIPPED
)
