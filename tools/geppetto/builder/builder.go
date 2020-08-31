package builder

import (
	"fmt"

	"ponglehub.co.uk/geppetto/services"

	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/geppetto/types"
)

type worker interface {
	buildNPM(repo types.Repo, signals chan<- buildSignal)
}

// Builder builds your application
type Builder struct {
	worker worker
}

// New create a new builder object
func New() *Builder {
	return &Builder{
		worker: &defaultWorker{
			npm: services.NewNpmService(),
		},
	}
}

// Build build your repos
func (b *Builder) Build(repos []types.Repo) error {
	state := NewBuildState()
	signals := make(chan buildSignal)

	for {
		for _, repo := range repos {
			ok, block := state.CanBuild(repo.Name, repo.DependsOn)
			if ok {
				logrus.Infof("Repo building: %s", repo.Name)
				state.Build(repo.Name)
				switch repo.RepoType {
				case types.Node:
					go b.worker.buildNPM(repo, signals)
				case types.Golang:
					fallthrough
				case types.Helm:
					logrus.Infof("Skipping build for %s, %s repos not implemented yet", repo.Name, repo.RepoType)
					state.Complete(repo.Name)
				default:
					return fmt.Errorf("Unknown repo type: %s", repo.RepoType)
				}
			}
			if block {
				logrus.Infof("Repo blocked: %s", repo.Name)
				state.Block(repo.Name)
			}
		}

		count := state.Count(BuildingState)
		if count == 0 {
			break
		}
		logrus.Debugf("Building %d repos", count)

		signal := <-signals
		if signal.err != nil {
			logrus.Errorf("Failed to build %s: %+v", signal.repo, signal.err)
			state.Error(signal.repo)
			continue
		}

		logrus.Infof("Finished building repo: %s", signal.repo)
		state.Complete(signal.repo)
	}

	return nil
}
