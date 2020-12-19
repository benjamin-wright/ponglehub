package builder

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/geppetto/types"
)

type worker interface {
	buildNPM(ctx context.Context, repo types.Repo, reinstall bool, signals chan<- signal)
	buildHelm(ctx context.Context, repo types.Repo, reinstall bool, signals chan<- signal)
	buildGolang(ctx context.Context, repo types.Repo, reinstall bool, signals chan<- signal)
	buildRust(ctx context.Context, repo types.Repo, reinstall bool, signals chan<- signal)
}

// Builder builds your application
type Builder struct {
	worker worker
}

// New create a new builder object
func New(chartRepo string) *Builder {
	return &Builder{
		worker: newDefaultWorker(chartRepo),
	}
}

// Build build your repos
func (b *Builder) Build(repos []types.Repo, inputs <-chan InputSignal) <-chan []types.RepoState {
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
					buildContext, cancelFunc := context.WithCancel(context.Background())
					reinstall := state.find(repo.Name).Start(buildContext, cancelFunc)

					switch repo.RepoType {
					case types.Node:
						go b.worker.buildNPM(buildContext, repo, reinstall, signals)
					case types.Helm:
						go b.worker.buildHelm(buildContext, repo, reinstall, signals)
					case types.Golang:
						go b.worker.buildGolang(buildContext, repo, reinstall, signals)
					case types.Rust:
						go b.worker.buildRust(buildContext, repo, reinstall, signals)
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
			case input := <-inputs:
				if input.Invalidate {
					state.invalidate(input.Repo, input.Reinstall)
					if state.find(input.Repo).Pending() {
						logrus.Infof("Unlocking repo build: %s", input.Repo)
					}
					continue
				}

				if input.Nuke {
					for _, repo := range repos {
						state.invalidate(repo.Name, true)
						if state.find(repo.Name).Pending() {
							logrus.Infof("Unlocking repo build: %s", repo.Name)
						}
					}
					continue
				}
			case signal := <-signals:
				if signal.err != nil {
					logrus.Errorf("Failed to build %s: %+v", signal.repo, signal.err)
					state.find(signal.repo).Error(signal.err)
					continue
				}

				if signal.cancelled {
					logrus.Infof("Ignoring signal for cancelled build of repo: '%s'", signal.repo)
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
				repo := state.find(signal.repo)
				repo.Complete()

				for _, r := range state.repos {
					for _, dep := range r.Repo().DependsOn {
						if dep == repo.Repo().Name {
							logrus.Debugf("Invalidating %s with dependency %s", r.Repo().Name, dep)
							state.invalidate(r.Repo().Name, true)
						}
					}
				}
			}
		}
	}()

	return progress
}
