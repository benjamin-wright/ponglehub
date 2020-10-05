package services

import (
	"context"
	"errors"
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
func (c *Commander) Run(ctx context.Context, workDir string, command string) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", command)
	cmd.Dir = workDir

	type result struct {
		value string
		err   error
	}

	resultChan := make(chan result)

	go func(results chan<- result) {
		out, err := cmd.CombinedOutput()
		if err != nil {
			results <- result{
				value: strings.TrimSpace(string(out)),
				err:   err,
			}
		}

		logrus.Debugf("Command `%s` output:\n%s", command, string(out))
		results <- result{
			value: strings.TrimSpace(string(out)),
			err:   err,
		}
	}(resultChan)

	select {
	case r := <-resultChan:
		return r.value, r.err
	case <-ctx.Done():
		cmd.Process.Kill()
		return "", errors.New("command cancelled")
	}
}
