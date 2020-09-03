package scanner

import (
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"ponglehub.co.uk/geppetto/types"
)

func (s *Scanner) WatchDir(repos []types.Repo) (<-chan types.Repo, <-chan error, <-chan bool) {
	triggers := make(chan types.Repo)
	errors := make(chan error)
	stopper := make(chan bool)

	stoppers := []chan bool{}

	for _, repo := range repos {
		repoStopper := make(chan bool)
		stoppers = append(stoppers, repoStopper)

		switch repo.RepoType {
		case types.Node:
			go s.watchNpm(repo, triggers, errors, repoStopper)
		default:

		}
	}

	go func() {
		<-stopper
		for _, s := range stoppers {
			s <- true
		}
	}()

	return triggers, errors, stopper
}

func (s *Scanner) watchNpm(repo types.Repo, triggers chan<- types.Repo, errors chan<- error, stopper <-chan bool) {
	watcher, _ := fsnotify.NewWatcher()
	defer watcher.Close()

	err := filepath.Walk(repo.Path, func(path string, fi os.FileInfo, err error) error {
		if fi.Mode().IsDir() {
			return watcher.Add(path)
		}

		return nil
	})

	if err != nil {
		errors <- err
	}

	go func() {
		for {
			select {
			// watch for events
			case <-watcher.Events:
				triggers <- repo

				// watch for errors
			case err := <-watcher.Errors:
				errors <- err
			}
		}
	}()

	<-stopper
}
