package builder

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/geppetto/types"
)

type worker interface {
	buildNPM(repo types.Repo, signals chan<- signal)
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
func (b *Builder) Build(repos []types.Repo, progress chan<- []types.RepoState) error {
	state := newBuildState(repos)
	signals := make(chan signal)

	progress <- state.repos

	for {
		for _, repo := range repos {
			ok, block := state.canBuild(repo.Name)

			if ok {
				logrus.Infof("Repo building: %s", repo.Name)
				state.find(repo.Name).Start()

				switch repo.RepoType {
				case types.Node:
					go b.worker.buildNPM(repo, signals)
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
		if count == 0 {
			break
		}
		logrus.Debugf("Building %d repos", count)

		signal := <-signals
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

	return nil
}
