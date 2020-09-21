package builder

import (
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/geppetto/services"
	"ponglehub.co.uk/geppetto/types"
)

type defaultWorker struct {
	npm services.NPM
}

func newDefaultWorker() *defaultWorker {
	return &defaultWorker{
		npm: services.NewNpmService(),
	}
}

func (w *defaultWorker) buildNPM(repo types.Repo, signals chan<- signal) {
	logrus.Debugf("Building NPM repo: %s", repo.Name)

	signals <- signal{repo: repo.Name, phase: "check"}
	currentSHA, err := w.npm.GetCurrentSHA(repo)
	if err != nil {
		signals <- signal{repo: repo.Name, err: err}
		return
	}

	latestSHA, err := w.npm.GetLatestSHA(repo)
	if err != nil {
		signals <- signal{repo: repo.Name, err: err}
		return
	}

	if currentSHA == latestSHA {
		signals <- signal{repo: repo.Name, skip: true}
		return
	}

	if repo.Reinstall {
		signals <- signal{repo: repo.Name, phase: "install"}
		if err := w.npm.Install(repo); err != nil {
			signals <- signal{repo: repo.Name, err: err}
			return
		}
	}

	signals <- signal{repo: repo.Name, phase: "lint"}
	if err := w.npm.Lint(repo); err != nil {
		signals <- signal{repo: repo.Name, err: err}
		return
	}

	signals <- signal{repo: repo.Name, phase: "test"}
	if err := w.npm.Test(repo); err != nil {
		signals <- signal{repo: repo.Name, err: err}
		return
	}

	signals <- signal{repo: repo.Name, phase: "publish"}
	if err := w.npm.Publish(repo); err != nil {
		signals <- signal{repo: repo.Name, err: err}
		return
	}

	signals <- signal{repo: repo.Name, finished: true}
}
