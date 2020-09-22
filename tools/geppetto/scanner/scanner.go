package scanner

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/geppetto/services"
	"ponglehub.co.uk/geppetto/types"
)

// Scanner a file-system scanner for finding code projects
type Scanner struct {
	io   services.IO
	npm  services.NPM
	helm services.Helm
}

// New creates a new scanner instance
func New() *Scanner {
	return &Scanner{
		io:   services.IO{},
		npm:  services.NewNpmService(),
		helm: services.NewHelmService(),
	}
}

// ScanDir finds code directories and returns a list of Repo objects representing them
func (s *Scanner) ScanDir(targetDir string) ([]types.Repo, error) {
	repos := []types.Repo{}

	err := s.io.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
		ignore := []string{"node_modules", ".git"}
		isIgnore := func(name string) bool {
			for _, i := range ignore {
				if name == i {
					return true
				}
			}
			return false
		}

		name := info.Name()

		if !info.IsDir() {
			return nil
		}

		if isIgnore(name) {
			return filepath.SkipDir
		}

		if s.io.FileExists(path + "/chart.yaml") {
			logrus.Infof("HELM: %s", path)
			repo, err := s.helm.GetRepo(path)
			if err != nil {
				return err
			}

			repos = append(repos, repo)
			return filepath.SkipDir
		}

		if s.io.FileExists(path + "/package.json") {
			logrus.Infof("NPM: %s", path)
			repo, err := s.npm.GetRepo(path)
			if err != nil {
				return err
			}

			repos = append(repos, repo)
			return filepath.SkipDir
		}

		if s.io.FileExists(path + "/go.mod") {
			logrus.Infof("GOLANG: %s", path)
			repo := types.Repo{
				Name:      name,
				Path:      path,
				RepoType:  types.Golang,
				DependsOn: []string{},
			}

			repos = append(repos, repo)
			return filepath.SkipDir
		}

		logrus.Debugf("- Unrecognised: %s", path)
		return nil
	})

	if err != nil {
		return nil, err
	}

	err = s.linkNPMRepos(repos)
	if err != nil {
		return nil, err
	}

	return repos, nil
}

func (s *Scanner) linkNPMRepos(repos []types.Repo) error {
	names := s.getNPMModuleNames(repos)
	for index, repo := range repos {
		if repo.RepoType != types.Node {
			continue
		}

		deps, err := s.npm.GetDependencyNames(repo)
		if err != nil {
			return err
		}

		logrus.Debugf("Dependencies for %s: %v", repo.Name, deps)

		for _, name := range names {
			for _, dep := range deps {
				if name == dep {
					repos[index].DependsOn = append(repos[index].DependsOn, name)
					continue
				}
			}
		}

	}

	return nil
}

func (s *Scanner) getNPMModuleNames(repos []types.Repo) []string {
	names := []string{}

	for _, repo := range repos {
		if repo.RepoType != types.Node {
			continue
		}

		names = append(names, repo.Name)
	}

	return names
}
