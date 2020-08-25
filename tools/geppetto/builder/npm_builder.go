package builder

import (
	"os/exec"

	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/geppetto/config"
)

type npmBuilder struct {
	basePath string
}

func (b npmBuilder) build(repo config.Repo, signal chan<- buildSignal) {
	logrus.Infof("Running NPM build for %s", repo.Name)
	err := b.runCommand(repo, "npm", "install", "--strict-ssl=false")
	if err == nil {
		err = b.runCommand(repo, "npm", "run", "lint")
	}
	signal <- buildSignal{
		repo: repo,
		err:  err,
	}
}

func (b npmBuilder) runCommand(repo config.Repo, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Dir = b.basePath + "/" + repo.Path

	out, err := cmd.CombinedOutput()
	if err != nil {
		logrus.Error(string(out))
		return err
	}

	logrus.Debugf("Command %s %v output:\n%s", command, args, string(out))
	return nil
}
