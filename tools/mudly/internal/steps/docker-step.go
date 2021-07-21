package steps

import (
	"bufio"
	"os/exec"

	"github.com/sirupsen/logrus"
)

type DockerStep struct {
	Name       string `yaml:"name"`
	Dockerfile string `yaml:"dockerfile"`
	Context    string `yaml:"context"`
	Tag        string `yaml:"tag"`
}

func (d DockerStep) Run(env map[string]string) bool {
	cmd := exec.Command("docker")

	if d.Dockerfile != "" {
		cmd.Args = append(cmd.Args, "-f", d.Dockerfile)
	}

	if d.Tag != "" {
		cmd.Args = append(cmd.Args, "-t", d.Tag)
	}

	if d.Context != "" {
		cmd.Args = append(cmd.Args, d.Context)
	} else {
		cmd.Args = append(cmd.Args, ".")
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logrus.Errorf("%s: Failed to get docker step stdout pipe: %+v", d.Name, err)
		return false
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		logrus.Errorf("%s: Failed to get docker step stderr pipe: %+v", d.Name, err)
		return false
	}

	err = cmd.Start()
	if err != nil {
		logrus.Errorf("%s: Docker step couldn't run: %+v", d.Name, err)
		return false
	}

	stdoutScanner := bufio.NewScanner(stdout)
	stdoutScanner.Split(bufio.ScanLines)
	for stdoutScanner.Scan() {
		logrus.Infof("%s: %s", d.Name, stdoutScanner.Text())
	}

	stderrScanner := bufio.NewScanner(stderr)
	stderrScanner.Split(bufio.ScanLines)
	for stderrScanner.Scan() {
		logrus.Warnf("%s: %s", d.Name, stderrScanner.Text())
	}

	err = cmd.Wait()
	if err != nil {
		logrus.Errorf("%s: Command failed: %+v", d.Name, err)
		return false
	}

	return true
}

func (d DockerStep) String() string {
	return d.Name
}
