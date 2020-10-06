package builder

import (
	"context"

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

func (w *defaultWorker) buildGolang(ctx context.Context, repo types.Repo, reinstall bool, signals chan<- signal) {
	logrus.Debugf("Building Golang repo: %s", repo.Name)

	if reinstall {
		signals <- signal{repo: repo.Name, phase: "tidy"}
		if err := w.golang.Tidy(ctx, repo); err != nil {
			signals <- makeErrorSignal(ctx, repo.Name, err)
		}

		signals <- signal{repo: repo.Name, phase: "install"}
		if err := w.golang.Install(ctx, repo); err != nil {
			signals <- makeErrorSignal(ctx, repo.Name, err)
			return
		}
	}

	signals <- signal{repo: repo.Name, phase: "test"}
	if err := w.golang.Test(ctx, repo); err != nil {
		signals <- makeErrorSignal(ctx, repo.Name, err)
		return
	}

	if w.golang.Buildable(repo) {
		signals <- signal{repo: repo.Name, phase: "build"}
		if err := w.golang.Build(ctx, repo); err != nil {
			signals <- makeErrorSignal(ctx, repo.Name, err)
			return
		}
	}

	signals <- signal{repo: repo.Name, finished: true}
}

func (w *defaultWorker) buildHelm(ctx context.Context, repo types.Repo, reinstall bool, signals chan<- signal) {
	logrus.Debugf("Building Helm repo: %s", repo.Name)

	if reinstall {
		signals <- signal{repo: repo.Name, phase: "install"}
		if err := w.helm.Install(ctx, repo); err != nil {
			signals <- makeErrorSignal(ctx, repo.Name, err)
			return
		}
	}

	signals <- signal{repo: repo.Name, phase: "lint"}
	if err := w.helm.Lint(ctx, repo); err != nil {
		signals <- makeErrorSignal(ctx, repo.Name, err)
		return
	}

	signals <- signal{repo: repo.Name, phase: "bump"}
	if err := w.helm.SetVersion(repo, ""); err != nil {
		signals <- makeErrorSignal(ctx, repo.Name, err)
		return
	}

	signals <- signal{repo: repo.Name, phase: "publish"}
	if err := w.helm.Publish(ctx, repo, w.chartRepo); err != nil {
		signals <- makeErrorSignal(ctx, repo.Name, err)
		return
	}

	signals <- signal{repo: repo.Name, finished: true}
}

func (w *defaultWorker) buildNPM(ctx context.Context, repo types.Repo, reinstall bool, signals chan<- signal) {
	logrus.Debugf("Building NPM repo: %s", repo.Name)

	if reinstall {
		signals <- signal{repo: repo.Name, phase: "install"}
		if err := w.npm.Install(ctx, repo); err != nil {
			signals <- makeErrorSignal(ctx, repo.Name, err)
			return
		}
	} else if !repo.Application {
		signals <- signal{repo: repo.Name, phase: "check"}
		currentSHA, err := w.npm.GetCurrentSHA(ctx, repo)
		if err != nil {
			signals <- makeErrorSignal(ctx, repo.Name, err)
			return
		}

		latestSHA, err := w.npm.GetLatestSHA(ctx, repo)
		if err != nil {
			signals <- makeErrorSignal(ctx, repo.Name, err)
			return
		}

		if currentSHA == latestSHA {
			signals <- signal{repo: repo.Name, skip: true}
			return
		}
	}

	signals <- signal{repo: repo.Name, phase: "lint"}
	if err := w.npm.Lint(ctx, repo); err != nil {
		signals <- makeErrorSignal(ctx, repo.Name, err)
		return
	}

	signals <- signal{repo: repo.Name, phase: "test"}
	if err := w.npm.Test(ctx, repo); err != nil {
		signals <- makeErrorSignal(ctx, repo.Name, err)
		return
	}

	if repo.Application {
		signals <- signal{repo: repo.Name, phase: "build"}
		if err := w.npm.Build(ctx, repo); err != nil {
			signals <- makeErrorSignal(ctx, repo.Name, err)
			return
		}
	} else {
		signals <- signal{repo: repo.Name, phase: "bump"}
		if err := w.npm.SetVersion(ctx, repo, ""); err != nil {
			signals <- makeErrorSignal(ctx, repo.Name, err)
			return
		}

		signals <- signal{repo: repo.Name, phase: "publish"}
		if err := w.npm.Publish(ctx, repo); err != nil {
			signals <- makeErrorSignal(ctx, repo.Name, err)
			return
		}

		signals <- signal{repo: repo.Name, phase: "clean"}
		if err := w.npm.CleanTempFiles(ctx, repo); err != nil {
			signals <- makeErrorSignal(ctx, repo.Name, err)
			return
		}
	}

	signals <- signal{repo: repo.Name, finished: true}
}
