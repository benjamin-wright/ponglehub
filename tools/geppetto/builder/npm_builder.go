package builder

import (
	"os/exec"

	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/geppetto/builder/commands"
	"ponglehub.co.uk/geppetto/config"
)

type npmBuilder struct {
	basePath string
}

func (b npmBuilder) build(repo config.Repo, signal chan<- buildSignal) {
	logrus.Infof("Running NPM build for %s", repo.Name)
	name, err := b.getName(repo)
	if err != nil {
		logrus.Infof("Build for %s failed", repo.Name)
		signal <- buildSignal{
			repo: repo,
			err:  err,
		}

		return
	}

	logrus.Infof("%s -> %s", repo.Name, name)

	for _, cmd := range []commands.Command{
		commands.MakeNpmCommand(b.basePath, &repo, "install", []string{"install", "--strict-ssl=false"}),
		commands.MakeNpmCommand(b.basePath, &repo, "lint", []string{"run", "lint"}),
		commands.MakeNpmCommand(b.basePath, &repo, "test", []string{"run", "test"}),
		commands.MakeNpmCommand(b.basePath, &repo, "publish", []string{"publish", "--strict-ssl=false"}),
	} {
		logrus.Infof(" - %s stage %s", repo.Name, cmd.Stage())
		err := cmd.Run()
		if err != nil {
			logrus.Infof("Build for %s failed", repo.Name)
			signal <- buildSignal{
				repo: repo,
				err:  err,
			}

			return
		}
	}

	logrus.Infof("Build for %s finished", repo.Name)
	signal <- buildSignal{
		repo: repo,
		err:  nil,
	}
}

func (b npmBuilder) getName(repo config.Repo) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", "cat package.json | jq '.name' -r")
	cmd.Dir = b.basePath + "/" + repo.Path

	out, err := cmd.CombinedOutput()
	if err != nil {
		logrus.Error(string(out))
		return "", err
	}

	return string(out), nil
}

// func (b npmBuilder) getLatestShasum(path string, name string) (string, error) {
// 	cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("npm view %s --strict-ssl=false --json | jq '.dist.shasum' -r"))
// 	cmd.Dir = b.basePath + "/" + path

// 	out, err := cmd.CombinedOutput()
// 	if err != nil {
// 		logrus.Error(string(out))
// 		return "", err
// 	}

// 	logrus.Debugf("Command npm %v output:\n%s", c.args, string(out))
// 	return nil
// 	return "npm view @pongle/eslint-config-ponglehub --strict-ssl=false --json | jq '.dist.shasum' -r"
// }

// func (b npmBuilder) getCurrentShasum(repo *config.Repo) string {
// 	return "npm publish --dry-run --json | jq '.shasum' -r"
// }
