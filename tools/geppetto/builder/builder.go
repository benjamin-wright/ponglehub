package builder

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/geppetto/types"
)

type worker interface {
	buildNPM(repo types.Repo, reinstall bool, signals chan<- signal)
}

// Builder builds your application
type Builder struct {
	worker worker
}

// New create a new builder object
func New() *Builder {
	return &Builder{
		worker: newDefaultWorker(),
	}
}

// Build build your repos
func (b *Builder) Build(repos []types.Repo, updates <-chan types.RepoUpdate) <-chan []types.RepoState {
	state := newBuildState(repos)
	signals := make(chan signal)
	progress := make(chan []types.RepoState, 5)

	progress <- state.repos

	go func() {
		for {
			for _, repo := range repos {
				ok, block := state.canBuild(repo.Name)

				if ok {
					logrus.Infof("Repo building: %s", repo.Name)
					reinstall := state.find(repo.Name).Start()

					switch repo.RepoType {
					case types.Node:
						go b.worker.buildNPM(repo, reinstall, signals)
					case types.Golang:
						fallthrough
					case types.Helm:
						logrus.Infof("Skipping build for %s, %s repos not implemented yet", repo.Name, repo.RepoType)
						state.find(repo.Name).Block()
					default:
						state.find(repo.Name).Error(fmt.Errorf("Unknown repo type: %s", repo.RepoType))
					}
				}

				if block {
					logrus.Infof("Repo blocked: %s", repo.Name)
					state.find(repo.Name).Block()
				}
			}

			progress <- state.repos

			count := state.numBuilding()
			if count > 0 {
				logrus.Debugf("Building %d repos", count)
			} else {
				progress <- nil
				logrus.Debugf("Waiting for updates...")
			}

			select {
			case update := <-updates:
				state.invalidate(update.Name, update.Install)
			case signal := <-signals:
				if signal.err != nil {
					logrus.Errorf("Failed to build %s: %+v", signal.repo, signal.err)
					state.find(signal.repo).Error(signal.err)
					continue
				}

				if signal.skip {
					logrus.Infof("Skipping repo: %s", signal.repo)
					state.find(signal.repo).Skip()
					continue
				}

				if signal.phase != "" {
					state.find(signal.repo).Progress(signal.phase)
					continue
				}

				logrus.Infof("Finished building repo: %s", signal.repo)
				state.find(signal.repo).Complete()
			}
		}

		close(progress)
	}()

	return progress
}
