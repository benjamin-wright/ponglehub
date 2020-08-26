package builder

import (
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/geppetto/builder/commands"
	"ponglehub.co.uk/geppetto/config"
	"ponglehub.co.uk/geppetto/services"
)

type npmBuilder struct {
	basePath string
}

func (b npmBuilder) build(repo config.Repo, signal chan<- buildSignal) {
	logrus.Infof("Running NPM build for %s", repo.Name)

	service, err := services.NewNpmRepo(b.basePath + "/" + repo.Path)
	if err != nil {
		logrus.Errorf("Build for %s failed: %+v", repo.Name, err)
	}

	for _, cmd := range []commands.Command{
		commands.MakeNpmCheckCommand(b.basePath, &repo),
		commands.CreateGeneric("install", func() (bool, error) { return false, service.Install() }),
		commands.CreateGeneric("lint", func() (bool, error) { return false, service.Lint() }),
		commands.CreateGeneric("test", func() (bool, error) { return false, service.Test() }),
		commands.CreateGeneric("publish", func() (bool, error) { return false, service.Publish() }),
	} {
		logrus.Infof(" - %s stage %s", repo.Name, cmd.Stage())
		skip, err := cmd.Run()

		if err != nil {
			logrus.Infof("Build for %s failed", repo.Name)
			signal <- buildSignal{repo: repo, err: err}
			return
		}

		if skip {
			logrus.Infof("Build for %s skipped", repo.Name)
			signal <- buildSignal{repo: repo}
			return
		}
	}

	logrus.Infof("Build for %s finished", repo.Name)
	signal <- buildSignal{repo: repo}
}
