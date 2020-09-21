package scanner

import (
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/geppetto/types"
)

func (s *Scanner) WatchDir(repos []types.Repo) (<-chan types.Repo, <-chan error) {
	triggers := make(chan types.Repo)
	errors := make(chan error)

	for _, repo := range repos {
		switch repo.RepoType {
		case types.Node:
			go s.watchNpm(repo, triggers, errors)
		default:
		}
	}

	return triggers, errors
}

func (s *Scanner) watchNpm(repo types.Repo, triggers chan<- types.Repo, errors chan<- error) {
	watcher, _ := fsnotify.NewWatcher()
	defer watcher.Close()

	err := filepath.Walk(repo.Path, func(path string, fi os.FileInfo, err error) error {
		if fi.Name() == "node_modules" {
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
			repo.Reinstall = filepath.Base(event.Name) == "package.json" || filepath.Base(event.Name) == "package-lock.json"
			triggers <- repo
			logrus.Infof("Sending trigger for %s for %s", repo.Name, filepath.Base(event.Name))

		// watch for errors
		case err := <-watcher.Errors:
			logrus.Infof("Error! %s", repo.Name)
			errors <- err
		}
	}
}
