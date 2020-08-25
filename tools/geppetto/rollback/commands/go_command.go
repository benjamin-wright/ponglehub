package commands

import (
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/geppetto/config"
)

// GoCommand represents a command to Go
type GoCommand struct {
	basePath string
	repo     config.Repo
}

// MakeGoCommand make a new Go rollback command object
func MakeGoCommand(basePath string, repo config.Repo) *GoCommand {
	return &GoCommand{
		basePath: basePath,
		repo:     repo,
	}
}

// Run run the Go command and return an error if it fails
func (e GoCommand) Run() error {
	logrus.Debugf("Running rollback on %s", e.repo.Name)
	return nil
}

// Name the name of the job
func (e GoCommand) Name() string {
	return e.repo.Name
}
