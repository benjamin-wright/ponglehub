package commands

import (
	"os/exec"

	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/geppetto/config"
)

// NpmCommand represents a command to NPM
type NpmCommand struct {
	basePath string
	repo     *config.Repo
	stage    string
	args     []string
}

// MakeNpmCommand make a new NPM command object
func MakeNpmCommand(basePath string, repo *config.Repo, stage string, args []string) *NpmCommand {
	return &NpmCommand{
		basePath: basePath,
		repo:     repo,
		stage:    stage,
		args:     args,
	}
}

// Run run the NPM command and return an error if it fails
func (e NpmCommand) Run() error {
	cmd := exec.Command("npm", e.args...)
	cmd.Dir = e.basePath + "/" + e.repo.Path

	out, err := cmd.CombinedOutput()
	if err != nil {
		logrus.Error(string(out))
		return err
	}

	logrus.Debugf("Command npm %v output:\n%s", e.args, string(out))
	return nil
}

// Stage get the name of the build stage
func (e NpmCommand) Stage() string {
	return e.stage
}
