package steps

import (
	"bufio"
	"fmt"
	"os/exec"

	"github.com/sirupsen/logrus"
)

type CommandStep struct {
	Name    string            `yaml:"name"`
	Command string            `yaml:"cmd"`
	Env     map[string]string `yaml:"env"`
}

func (c CommandStep) Run(env map[string]string) bool {
	cmd := exec.Command("/bin/bash", "-c", c.Command)

	for key, value := range c.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logrus.Errorf("%s: Failed to get command stdout pipe: %+v", c.Name, err)
		return false
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		logrus.Errorf("%s: Failed to get command stderr pipe: %+v", c.Name, err)
		return false
	}

	err = cmd.Start()
	if err != nil {
		logrus.Errorf("%s: Command couldn't run: %+v", c.Name, err)
		return false
	}

	stdoutScanner := bufio.NewScanner(stdout)
	stdoutScanner.Split(bufio.ScanLines)
	for stdoutScanner.Scan() {
		logrus.Infof("%s: %s", c.Name, stdoutScanner.Text())
	}

	stderrScanner := bufio.NewScanner(stderr)
	stderrScanner.Split(bufio.ScanLines)
	for stderrScanner.Scan() {
		logrus.Warnf("%s: %s", c.Name, stderrScanner.Text())
	}

	err = cmd.Wait()
	if err != nil {
		logrus.Errorf("%s: Command failed: %+v", c.Name, err)
		return false
	}

	return true
}

func (c CommandStep) String() string {
	return c.Name
}
