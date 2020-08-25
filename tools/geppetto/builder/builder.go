package builder

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/geppetto/config"
)

type buildAgent interface {
	build(repo config.Repo, signal chan<- buildSignal)
}

// Builder coordinate builds
type Builder struct {
	plan     planner
	cfg      *config.Config
	signals  chan buildSignal
	npmAgent buildAgent
	goAgent  buildAgent
}

type stepKind string

type buildSignal struct {
	repo config.Repo
	err  error
}

// FromConfig create a new builder instance from the given config
func FromConfig(cfg *config.Config) (*Builder, error) {
	builder := Builder{
		plan: planner{
			built: []string{},
		},
		cfg:      cfg,
		signals:  make(chan buildSignal),
		npmAgent: npmBuilder{basePath: cfg.BasePath},
		goAgent:  goBuilder{basePath: cfg.BasePath},
	}

	if collisions, ok := builder.hasCircularDependencies(); !ok {
		return nil, fmt.Errorf("Error creating builder from config, detected circular dependency affecting %+v", collisions)
	}

	return &builder, nil
}

func (builder *Builder) hasCircularDependencies() (collisions []string, ok bool) {
	builder.plan.reset()
	building := true

	for building {
		building = false
		builds := []string{}

		for _, r := range builder.cfg.Repos {
			if builder.plan.isBuilt(r.Name) {
				continue
			}

			if builder.plan.areBuilt(r.DependsOn) {
				building = true
				builds = append(builds, r.Name)
			}
		}

		for _, build := range builds {
			builder.plan.complete(build)
		}
	}

	if len(builder.plan.built) != len(builder.cfg.Repos) {
		pending := []string{}

		for _, repo := range builder.cfg.Repos {
			if !builder.plan.isBuilt(repo.Name) {
				pending = append(pending, repo.Name)
			}
		}

		return pending, false
	}

	return []string{}, true
}

// Build run a full build of everything
func (builder *Builder) Build() error {
	builder.plan.reset()
	running := true

	for running {
		running = false

		for _, repo := range builder.cfg.Repos {
			if builder.plan.canBuild(repo) {
				logrus.Debugf("Can build %s", repo.Name)
				running = true
				builder.plan.run(repo.Name)

				switch repo.RepoType {
				case config.Node:
					go builder.npmAgent.build(repo, builder.signals)
				case config.Go:
					go builder.goAgent.build(repo, builder.signals)
				}
			}
		}

		if running {
			s := <-builder.signals

			builder.plan.complete(s.repo.Name)

			if s.err != nil {
				return s.err
			}
		}
	}

	return nil
}
