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
	io     services.IO
	npm    services.NPM
	helm   services.Helm
	golang services.Golang
	rust   services.Rust
}

// New creates a new scanner instance
func New() *Scanner {
	return &Scanner{
		io:     services.IO{},
		npm:    services.NewNpmService(),
		helm:   services.NewHelmService(),
		golang: services.NewGolangService(),
		rust:   services.NewRustService(),
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
			repo, err := s.golang.GetRepo(path)
			if err != nil {
				return err
			}
			logrus.Infof("repo: %+v", repo)

			repos = append(repos, repo)
			return filepath.SkipDir
		}

		if s.io.FileExists(path + "/Cargo.toml") {
			logrus.Infof("RUST: %s", path)
			repo, err := s.rust.GetRepo(path)
			if err != nil {
				return err
			}
			logrus.Infof("repo: %+v", repo)

			repos = append(repos, repo)
			return filepath.SkipDir
		}

		logrus.Debugf("- Unrecognised: %s", path)
		return nil
	})

	if err != nil {
		return nil, err
	}

	logrus.Info("Linking npm repos")
	err = s.linkNPMRepos(repos)
	if err != nil {
		return nil, err
	}

	logrus.Info("Linking helm repos")
	err = s.linkHelmRepos(repos)
	if err != nil {
		return nil, err
	}

	logrus.Info("Linking golang repos")
	err = s.linkGolangRepos(repos)
	if err != nil {
		return nil, err
	}

	return repos, nil
}

func (s *Scanner) linkNPMRepos(repos []types.Repo) error {
	names := s.getModuleNames(repos, types.Node)
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

func (s *Scanner) linkHelmRepos(repos []types.Repo) error {
	names := s.getModuleNames(repos, types.Helm)
	for index, repo := range repos {
		if repo.RepoType != types.Helm {
			continue
		}

		deps, err := s.helm.GetDependencyNames(repo)
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

func (s *Scanner) linkGolangRepos(repos []types.Repo) error {
	names := s.getModuleNames(repos, types.Golang)
	for index, repo := range repos {
		if repo.RepoType != types.Golang {
			continue
		}

		deps, err := s.golang.GetDependencyNames(repo)
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

func (s *Scanner) getModuleNames(repos []types.Repo, repoType types.RepoType) []string {
	names := []string{}

	for _, repo := range repos {
		if repo.RepoType != repoType {
			continue
		}

		names = append(names, repo.Name)
	}

	return names
}
