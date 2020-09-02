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

func stateToProgress(repos []types.Repo, state buildState) []types.RepoStatus {
	statuses := []types.RepoStatus{}

	for _, repo := range repos {
		repoState := state.GetState(repo.Name)
		repoPhase := state.GetPhase(repo.Name)

		statuses = append(statuses, types.RepoStatus{
			Repo:     repo,
			Blocked:  repoState == blockedState,
			Building: repoState == buildingState,
			Built:    repoState == builtState,
			Skipped:  repoState == skippedState,
			Error:    repoState == erroredState,
			Phase:    repoPhase,
		})
	}

	return statuses
}

// Build build your repos
func (b *Builder) Build(repos []types.Repo, progress chan<- []types.RepoStatus) error {
	state := newBuildState()
	signals := make(chan buildSignal)

	progress <- stateToProgress(repos, state)

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
					state.Error(repo.Name)
					return fmt.Errorf("Unknown repo type: %s", repo.RepoType)
				}
			}
			if block {
				logrus.Infof("Repo blocked: %s", repo.Name)
				state.Block(repo.Name)
			}
		}

		progress <- stateToProgress(repos, state)

		count := state.Count(buildingState)
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

		if signal.skip {
			logrus.Infof("Skipping repo: %s", signal.repo)
			state.Skip(signal.repo)
			continue
		}

		if signal.phase != "" {
			state.Progress(signal.repo, signal.phase)
			continue
		}

		logrus.Infof("Finished building repo: %s", signal.repo)
		state.Complete(signal.repo)
	}

	return nil
}
