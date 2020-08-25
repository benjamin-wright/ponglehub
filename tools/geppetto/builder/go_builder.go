package builder

import (
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/geppetto/config"
)

type goBuilder struct {
	basePath string
}

func (b goBuilder) build(repo config.Repo, signal chan<- buildSignal) {
	logrus.Infof("Running GO build for %s", repo.Name)
	signal <- buildSignal{
		repo: repo,
		err:  nil,
	}
}
