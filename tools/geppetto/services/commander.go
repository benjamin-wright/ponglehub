package services

import (
	"context"
	"errors"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

type commander interface {
	Run(ctx context.Context, workDir string, command string) (string, error)
}

// Commander runs shell commands
type Commander struct{}

// Run execute the command in the given directory and return the combined console output
func (c *Commander) Run(ctx context.Context, workDir string, command string) (string, error) {
	select {
	case <-ctx.Done():
		return "", errors.New("Canceled")
	default:
	}

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
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
		return "", errors.New("command cancelled")
	}
}
