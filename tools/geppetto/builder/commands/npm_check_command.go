package commands

import (
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/geppetto/config"
)

// NpmCheckCommand represents a command to check the version and published status of the module
type NpmCheckCommand struct {
	basePath string
	repo     *config.Repo
}

// MakeNpmCheckCommand make a new NPM check command object
func MakeNpmCheckCommand(basePath string, repo *config.Repo) *NpmCheckCommand {
	return &NpmCheckCommand{
		basePath: basePath,
		repo:     repo,
	}
}

// Run run the NPM command and return an error if it fails
func (e NpmCheckCommand) Run() (bool, error) {
	path := e.basePath + "/" + e.repo.Path

	oldSHA, err := runNpmCommand(path, "npm view --strict-ssl=false --json | jq '.dist.shasum' -r")
	if err != nil {
		return false, err
	}

	newSHA, err := runNpmCommand(path, "npm publish --dry-run --json | jq '.shasum' -r")
	if err != nil {
		return false, err
	}

	logrus.Debugf("Repo check for %s: %s -> %s", e.repo.Name, oldSHA, newSHA)

	if oldSHA == newSHA {
		return true, nil
	}

	return false, nil
}

// Stage get the name of the build stage
func (e NpmCheckCommand) Stage() string {
	return "check"
}

func runNpmCommand(path string, args string) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", args)
	cmd.Dir = path

	out, err := cmd.CombinedOutput()
	if err != nil {
		logrus.Error(string(out))
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}
