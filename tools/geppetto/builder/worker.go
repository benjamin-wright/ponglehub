package builder

import (
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/geppetto/services"
	"ponglehub.co.uk/geppetto/types"
)

type defaultWorker struct {
	npm services.NPM
}

func (w *defaultWorker) buildNPM(repo types.Repo, signals chan<- buildSignal) {
	logrus.Debugf("Building NPM repo: %s", repo.Name)

	currentSHA, err := w.npm.GetCurrentSHA(repo)
	if err != nil {
		signals <- buildSignal{repo: repo.Name, err: err}
		return
	}

	latestSHA, err := w.npm.GetLatestSHA(repo)
	if err != nil {
		signals <- buildSignal{repo: repo.Name, err: err}
		return
	}

	if currentSHA == latestSHA {
		signals <- buildSignal{repo: repo.Name, skip: true}
		return
	}

	if err := w.npm.Install(repo); err != nil {
		signals <- buildSignal{repo: repo.Name, err: err}
		return
	}

	if err := w.npm.Lint(repo); err != nil {
		signals <- buildSignal{repo: repo.Name, err: err}
		return
	}

	if err := w.npm.Test(repo); err != nil {
		signals <- buildSignal{repo: repo.Name, err: err}
		return
	}

	if err := w.npm.Publish(repo); err != nil {
		signals <- buildSignal{repo: repo.Name, err: err}
		return
	}

	signals <- buildSignal{repo: repo.Name}
}
