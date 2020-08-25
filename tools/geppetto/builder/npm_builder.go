package builder

import (
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/geppetto/builder/commands"
	"ponglehub.co.uk/geppetto/config"
)

type npmBuilder struct {
	basePath string
}

func (b npmBuilder) build(repo config.Repo, signal chan<- buildSignal) {
	logrus.Infof("Running NPM build for %s", repo.Name)

	for _, cmd := range []commands.Command{
		commands.MakeNpmCheckCommand(b.basePath, &repo),
		commands.MakeNpmCommand(b.basePath, &repo, "install", []string{"install", "--strict-ssl=false"}),
		commands.MakeNpmCommand(b.basePath, &repo, "lint", []string{"run", "lint"}),
		commands.MakeNpmCommand(b.basePath, &repo, "test", []string{"run", "test"}),
		commands.MakeNpmCommand(b.basePath, &repo, "publish", []string{"publish", "--strict-ssl=false"}),
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
