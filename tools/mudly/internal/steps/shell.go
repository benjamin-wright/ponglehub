package steps

import (
	"bufio"
	"fmt"
	"os/exec"
	"path"

	"github.com/sirupsen/logrus"
)

type shellCommand struct {
	dir      string
	artefact string
	step     string
	command  string
	args     []string
	env      map[string]string
	test     bool
}

func runShellCommand(command *shellCommand) bool {
	cmd := exec.Command(command.command, command.args...)
	cmd.Dir = path.Clean(command.dir)

	for key, value := range command.env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logrus.Errorf("%s[%s]: Failed to get command stdout pipe: %+v", command.artefact, command.step, err)
		return false
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		logrus.Errorf("%s[%s]: Failed to get command stderr pipe: %+v", command.artefact, command.step, err)
		return false
	}

	err = cmd.Start()
	if err != nil {
		logrus.Errorf("%s[%s]: Command couldn't run: %+v", command.artefact, command.step, err)
		return false
	}

	stdoutScanner := bufio.NewScanner(stdout)
	stdoutScanner.Split(bufio.ScanLines)
	for stdoutScanner.Scan() {
		logrus.Infof("%s[%s]: %s", command.artefact, command.step, stdoutScanner.Text())
	}

	stderrScanner := bufio.NewScanner(stderr)
	stderrScanner.Split(bufio.ScanLines)
	for stderrScanner.Scan() {
		logrus.Warnf("%s[%s]: %s", command.artefact, command.step, stderrScanner.Text())
	}

	err = cmd.Wait()
	if err != nil {
		if !command.test {
			logrus.Errorf("%s[%s]: Command failed: %+v", command.artefact, command.step, err)
		}
		return false
	}

	return true
}
