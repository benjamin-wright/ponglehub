package builder

import (
	"os/exec"

	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/geppetto/config"
)

type npmBuilder struct{}

func (b npmBuilder) build(repo config.Repo, signal chan<- buildSignal) {
	logrus.Infof("Running NPM build for %s", repo.Name)
	signal <- buildSignal{
		repo: repo,
		err:  nil,
	}
}

func (b npmBuilder) lint(repo config.Repo) error {
	cmd := exec.Command("npm lint")
}
