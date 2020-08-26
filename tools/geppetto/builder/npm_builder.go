package builder

import (
	"github.com/sirupsen/logrus"
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

	for _, s := range []struct {
		name string
		run  func() (bool, error)
	}{
		{
			name: "check",
			run: func() (bool, error) {
				oldSHA, err := service.GetCurrentSHA()
				if err != nil {
					return false, err
				}

				newSHA, err := service.GetLatestSHA()
				if err != nil {
					return false, err
				}

				return oldSHA == newSHA, nil
			},
		},
		{name: "install", run: func() (bool, error) { return false, service.Install() }},
		{name: "lint", run: func() (bool, error) { return false, service.Lint() }},
		{name: "test", run: func() (bool, error) { return false, service.Test() }},
		{name: "publish", run: func() (bool, error) { return false, service.Publish() }},
	} {
		logrus.Infof(" - %s stage %s", repo.Name, s.name)
		skip, err := s.run()

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
