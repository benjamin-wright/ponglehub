package scanner

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/geppetto/types"
)

// WatchDir watches for file changes in the supplied repos
func (s *Scanner) WatchDir(repos []types.Repo) (<-chan types.RepoUpdate, <-chan error) {
	triggers := make(chan types.RepoUpdate)
	errors := make(chan error)

	for _, repo := range repos {
		switch repo.RepoType {
		case types.Node:
			go s.watchNpm(repo, triggers, errors)
		case types.Helm:
			go s.watchHelm(repo, triggers, errors)
		case types.Golang:
			go s.watchGo(repo, triggers, errors)
		case types.Rust:
			go s.watchRust(repo, triggers, errors)
		default:
		}
	}

	return triggers, errors
}

func (s *Scanner) watchRust(repo types.Repo, triggers chan<- types.RepoUpdate, errors chan<- error) {
	watcher, _ := fsnotify.NewWatcher()
	defer watcher.Close()

	err := filepath.Walk(repo.Path, func(path string, fi os.FileInfo, err error) error {
		if fi.Name() == "target" {
			logrus.Infof("Skipping %s", fi.Name())
			return filepath.SkipDir
		}

		if fi.Name() == "build" {
			logrus.Infof("Skipping %s", fi.Name())
			return filepath.SkipDir
		}

		if fi.Mode().IsDir() {
			logrus.Infof("Monitoring %s", fi.Name())
			return watcher.Add(path)
		}

		return nil
	})

	if err != nil {
		errors <- err
	}

	for {
		select {
		// watch for events
		case event := <-watcher.Events:
			baseName := filepath.Base(event.Name)
			if strings.HasPrefix(baseName, "target") && !strings.HasSuffix(baseName, ".rs") {
				continue
			}

			if strings.HasSuffix(baseName, ".lock") {
				continue
			}

			triggers <- types.RepoUpdate{
				Name:    repo.Name,
				Path:    baseName,
				Install: baseName == "Cargo.toml",
			}

			logrus.Infof("Sending trigger for %s for %s", repo.Name, baseName)

		// watch for errors
		case err := <-watcher.Errors:
			logrus.Infof("Error! %s", repo.Name)
			errors <- err
		}
	}
}

func (s *Scanner) watchGo(repo types.Repo, triggers chan<- types.RepoUpdate, errors chan<- error) {
	watcher, _ := fsnotify.NewWatcher()
	defer watcher.Close()

	err := filepath.Walk(repo.Path, func(path string, fi os.FileInfo, err error) error {
		if fi.Name() == "build" {
			logrus.Infof("Skipping %s", fi.Name())
			return filepath.SkipDir
		}

		if fi.Mode().IsDir() {
			logrus.Infof("Monitoring %s", fi.Name())
			return watcher.Add(path)
		}

		return nil
	})

	if err != nil {
		errors <- err
	}

	for {
		select {
		// watch for events
		case event := <-watcher.Events:
			if strings.HasSuffix(filepath.Base(event.Name), ".log") || filepath.Base(event.Name) == "go.sum" {
				continue
			}

			triggers <- types.RepoUpdate{
				Name:    repo.Name,
				Path:    filepath.Base(event.Name),
				Install: filepath.Base(event.Name) == "go.mod",
			}

			logrus.Infof("Sending trigger for %s for %s", repo.Name, filepath.Base(event.Name))

		// watch for errors
		case err := <-watcher.Errors:
			logrus.Infof("Error! %s", repo.Name)
			errors <- err
		}
	}
}

func (s *Scanner) watchHelm(repo types.Repo, triggers chan<- types.RepoUpdate, errors chan<- error) {
	watcher, _ := fsnotify.NewWatcher()
	defer watcher.Close()

	err := filepath.Walk(repo.Path, func(path string, fi os.FileInfo, err error) error {
		if fi.Name() == "tmpcharts" {
			logrus.Infof("Skipping %s", fi.Name())
			return filepath.SkipDir
		}

		if fi.Mode().IsDir() {
			logrus.Infof("Monitoring %s", fi.Name())
			return watcher.Add(path)
		}

		return nil
	})

	if err != nil {
		errors <- err
	}

	for {
		select {
		// watch for events
		case event := <-watcher.Events:
			baseFile := filepath.Base(event.Name)

			if strings.HasSuffix(baseFile, ".tgz") || baseFile == "Chart.lock" {
				continue
			}

			if baseFile == "tmpcharts" || baseFile == "charts" {
				continue
			}

			triggers <- types.RepoUpdate{
				Name:    repo.Name,
				Path:    baseFile,
				Install: baseFile == "Chart.yaml",
			}

			logrus.Infof("Sending trigger for %s for %s", repo.Name, baseFile)

		// watch for errors
		case err := <-watcher.Errors:
			logrus.Infof("Error! %s", repo.Name)
			errors <- err
		}
	}
}

func (s *Scanner) watchNpm(repo types.Repo, triggers chan<- types.RepoUpdate, errors chan<- error) {
	watcher, _ := fsnotify.NewWatcher()
	defer watcher.Close()

	err := filepath.Walk(repo.Path, func(path string, fi os.FileInfo, err error) error {
		if fi.Name() == "node_modules" {
			logrus.Infof("Skipping %s", fi.Name())
			return filepath.SkipDir
		}

		if fi.Name() == "dist" {
			logrus.Infof("Skipping %s", fi.Name())
			return filepath.SkipDir
		}

		if fi.Mode().IsDir() {
			logrus.Infof("Monitoring %s", fi.Name())
			return watcher.Add(path)
		}

		return nil
	})

	if err != nil {
		errors <- err
	}

	for {
		select {
		// watch for events
		case event := <-watcher.Events:
			if strings.Contains(filepath.Base(event.Name), "package-lock.json") {
				continue
			}

			triggers <- types.RepoUpdate{
				Name:    repo.Name,
				Path:    filepath.Base(event.Name),
				Install: filepath.Base(event.Name) == "package.json",
			}

			logrus.Infof("Sending trigger for %s for %s", repo.Name, filepath.Base(event.Name))

		// watch for errors
		case err := <-watcher.Errors:
			logrus.Infof("Error! %s", repo.Name)
			errors <- err
		}
	}
}
