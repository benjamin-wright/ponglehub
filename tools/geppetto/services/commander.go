package services

import (
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

type commander interface {
	Run(workDir string, command string) (string, error)
}

// Commander runs shell commands
type Commander struct{}

// Run execute the command in the given directory and return the combined console output
func (c *Commander) Run(workDir string, command string) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", command)
	cmd.Dir = workDir

	out, err := cmd.CombinedOutput()
	if err != nil {
		return strings.TrimSpace(string(out)), err
	}

	logrus.Debugf("Command `%s` output:\n%s", command, string(out))
	return strings.TrimSpace(string(out)), nil
}
