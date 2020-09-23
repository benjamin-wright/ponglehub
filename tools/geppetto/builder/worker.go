package builder

import (
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/geppetto/services"
	"ponglehub.co.uk/geppetto/types"
)

type defaultWorker struct {
	npm       services.NPM
	helm      services.Helm
	golang    services.Golang
	chartRepo string
}

func newDefaultWorker(chartRepo string) *defaultWorker {
	return &defaultWorker{
		npm:       services.NewNpmService(),
		helm:      services.NewHelmService(),
		golang:    services.NewGolangService(),
		chartRepo: chartRepo,
	}
}

func (w *defaultWorker) buildGolang(repo types.Repo, reinstall bool, signals chan<- signal) {
	logrus.Debugf("Building Golang repo: %s", repo.Name)

	if reinstall {
		signals <- signal{repo: repo.Name, phase: "tidy"}
		if err := w.golang.Tidy(repo); err != nil {
			signals <- signal{repo: repo.Name, err: err}
			return
		}

		signals <- signal{repo: repo.Name, phase: "install"}
		if err := w.golang.Install(repo); err != nil {
			signals <- signal{repo: repo.Name, err: err}
			return
		}
	}

	signals <- signal{repo: repo.Name, finished: true}
}

func (w *defaultWorker) buildHelm(repo types.Repo, reinstall bool, signals chan<- signal) {
	logrus.Debugf("Building Helm repo: %s", repo.Name)

	if reinstall {
		signals <- signal{repo: repo.Name, phase: "install"}
		if err := w.helm.Install(repo); err != nil {
			signals <- signal{repo: repo.Name, err: err}
			return
		}
	}

	signals <- signal{repo: repo.Name, phase: "lint"}
	if err := w.helm.Lint(repo); err != nil {
		signals <- signal{repo: repo.Name, err: err}
		return
	}

	signals <- signal{repo: repo.Name, phase: "publish"}
	if err := w.helm.Publish(repo, w.chartRepo); err != nil {
		signals <- signal{repo: repo.Name, err: err}
		return
	}

	signals <- signal{repo: repo.Name, finished: true}
}

func (w *defaultWorker) buildNPM(repo types.Repo, reinstall bool, signals chan<- signal) {
	logrus.Debugf("Building NPM repo: %s", repo.Name)

	if reinstall {
		signals <- signal{repo: repo.Name, phase: "install"}
		if err := w.npm.Install(repo); err != nil {
			signals <- signal{repo: repo.Name, err: err}
			return
		}
	} else {
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

	signals <- signal{repo: repo.Name, phase: "bump"}
	if err := w.npm.SetVersion(repo, ""); err != nil {
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
