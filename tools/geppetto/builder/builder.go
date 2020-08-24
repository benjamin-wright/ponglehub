package builder

import (
	"fmt"

	"ponglehub.co.uk/geppetto/config"
)

// Builder coordinate builds
type Builder struct {
	plan    planner
	cfg     *config.Config
	signals chan buildSignal
}

type buildSignal struct {
	repo config.Repo
}

// FromConfig create a new builder instance from the given config
func FromConfig(cfg *config.Config) (*Builder, error) {
	builder := Builder{
		plan: planner{
			built: []string{},
		},
		cfg:     cfg,
		signals: make(chan buildSignal),
	}

	if collisions, ok := builder.hasCircularDependencies(); !ok {
		return nil, fmt.Errorf("Error creating builder from config, detected circular dependency affecting %+v", collisions)
	}

	return &builder, nil
}

func (builder Builder) hasCircularDependencies() (collisions []string, ok bool) {
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
