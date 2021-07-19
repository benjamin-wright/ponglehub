package steps

type runner interface {
	runCommand(cmd string) bool
}

type commandRunner struct{}

// func (c *commandRunner) runCommand(cmd string) bool {

// }
