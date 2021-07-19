package solver

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"ponglehub.co.uk/tools/mudly/internal/config"
	"ponglehub.co.uk/tools/mudly/internal/target"
)

type getArtefactResult struct {
	Config   string
	Artefact string
	Pipeline string
}

func TestGetArtefact(t *testing.T) {
	for _, test := range []struct {
		Name     string
		Target   target.Target
		Configs  []config.Config
		Expected *getArtefactResult
	}{
		{
			Name:   "simple",
			Target: target.Target{Dir: ".", Artefact: "test-artefact"},
			Configs: []config.Config{
				{
					Path: ".",
					Artefacts: []config.Artefact{
						{
							Name: "test-artefact",
						},
					},
				},
			},
			Expected: &getArtefactResult{Config: ".", Artefact: "test-artefact"},
		},
		{
			Name:   "picks the right artefact",
			Target: target.Target{Dir: ".", Artefact: "other"},
			Configs: []config.Config{
				{
					Path: ".",
					Artefacts: []config.Artefact{
						{
							Name: "test-artefact",
						},
						{
							Name: "other",
						},
					},
				},
			},
			Expected: &getArtefactResult{Config: ".", Artefact: "other"},
		},
		{
			Name:   "picks the right artefact reverse",
			Target: target.Target{Dir: ".", Artefact: "other"},
			Configs: []config.Config{
				{
					Path: ".",
					Artefacts: []config.Artefact{
						{
							Name: "other",
						},
						{
							Name: "test-artefact",
						},
					},
				},
			},
			Expected: &getArtefactResult{Config: ".", Artefact: "other"},
		},
		{
			Name:   "picks the right config",
			Target: target.Target{Dir: "./subdir", Artefact: "test-artefact"},
			Configs: []config.Config{
				{
					Path: ".",
					Artefacts: []config.Artefact{
						{
							Name:     "test-artefact",
							Pipeline: config.Pipeline{Name: "firstConfig"},
						},
					},
				},
				{
					Path: "subdir",
					Artefacts: []config.Artefact{
						{
							Name:     "test-artefact",
							Pipeline: config.Pipeline{Name: "secondConfig"},
						},
					},
				},
			},
			Expected: &getArtefactResult{Config: "subdir", Artefact: "test-artefact", Pipeline: "secondConfig"},
		},
		{
			Name:   "picks the right config reversed",
			Target: target.Target{Dir: "./subdir", Artefact: "test-artefact"},
			Configs: []config.Config{
				{
					Path: "subdir",
					Artefacts: []config.Artefact{
						{
							Name:     "test-artefact",
							Pipeline: config.Pipeline{Name: "firstConfig"},
						},
					},
				},
				{
					Path: ".",
					Artefacts: []config.Artefact{
						{
							Name:     "test-artefact",
							Pipeline: config.Pipeline{Name: "secondConfig"},
						},
					},
				},
			},
			Expected: &getArtefactResult{Config: "subdir", Artefact: "test-artefact", Pipeline: "firstConfig"},
		},
	} {
		t.Run(test.Name, func(u *testing.T) {
			cfg, artefact, err := getArtefact(test.Target, test.Configs)

			assert.NoError(u, err, "didn't expect an error")

			if test.Expected != nil {
				if cfg != nil && artefact != nil {
					assert.Equal(u, test.Expected, &getArtefactResult{Config: cfg.Path, Artefact: artefact.Name, Pipeline: artefact.Pipeline.Name})
				} else {
					assert.Fail(u, "expected a config and artefact", "%+v, %+v", cfg, artefact)
				}
			}
		})
	}
}

func TestCollectDependencies(t *testing.T) {
	for _, test := range []struct {
		Name     string
		Targets  []target.Target
		Configs  []config.Config
		Expected []Link
	}{
		{
			Name: "should get nothing from nothing",
		},
		{
			Name: "should find local links",
			Targets: []target.Target{
				{Dir: ".", Artefact: "artefact-1"},
			},
			Configs: []config.Config{
				{
					Path: ".",
					Artefacts: []config.Artefact{
						{
							Name: "artefact-1",
							Dependencies: []target.Target{
								{Dir: ".", Artefact: "artefact-2"},
							},
						},
						{
							Name: "artefact-2",
						},
					},
				},
			},
			Expected: []Link{
				{
					Target: target.Target{Dir: ".", Artefact: "artefact-2"},
					Source: target.Target{Dir: ".", Artefact: "artefact-1"},
				},
			},
		},
		{
			Name: "should find remote links",
			Targets: []target.Target{
				{Dir: "subdir1", Artefact: "artefact-1"},
			},
			Configs: []config.Config{
				{
					Path: "subdir1",
					Artefacts: []config.Artefact{
						{
							Name: "artefact-1",
							Dependencies: []target.Target{
								{Dir: "../subdir2", Artefact: "artefact-2"},
							},
						},
						{
							Name: "artefact-2",
						},
					},
				},
				{
					Path: "subdir2",
					Artefacts: []config.Artefact{
						{
							Name: "artefact-2",
						},
					},
				},
			},
			Expected: []Link{
				{
					Target: target.Target{Dir: "subdir2", Artefact: "artefact-2"},
					Source: target.Target{Dir: "subdir1", Artefact: "artefact-1"},
				},
			},
		},
		{
			Name: "should find chained dependency links",
			Targets: []target.Target{
				{Dir: "subdir1", Artefact: "artefact-1"},
			},
			Configs: []config.Config{
				{
					Path: "subdir1",
					Artefacts: []config.Artefact{
						{
							Name: "artefact-1",
							Dependencies: []target.Target{
								{Dir: "../subdir2", Artefact: "artefact-2"},
							},
						},
						{
							Name: "artefact-2",
						},
					},
				},
				{
					Path: "subdir2",
					Artefacts: []config.Artefact{
						{
							Name: "artefact-2",
							Dependencies: []target.Target{
								{Dir: "../subdir1", Artefact: "artefact-2"},
							},
						},
					},
				},
			},
			Expected: []Link{
				{
					Target: target.Target{Dir: "subdir2", Artefact: "artefact-2"},
					Source: target.Target{Dir: "subdir1", Artefact: "artefact-1"},
				},
				{
					Target: target.Target{Dir: "subdir1", Artefact: "artefact-2"},
					Source: target.Target{Dir: "subdir2", Artefact: "artefact-2"},
				},
			},
		},
	} {
		t.Run(test.Name, func(u *testing.T) {
			links, err := collectDependencies(test.Targets, test.Configs)

			assert.NoError(u, err, "didn't expect an error")

			if test.Expected != nil {
				if links != nil {
					assert.Equal(u, test.Expected, links)
				} else {
					assert.Fail(u, "expected a list of links")
				}
			}
		})
	}
}

func TestGetDedupedTargets(t *testing.T) {
	for _, test := range []struct {
		Name     string
		Targets  []target.Target
		Links    []Link
		Expected []target.Target
	}{
		{
			Name: "should get nothing from nothing",
		},
		{
			Name: "should return non-duplicated targets",
			Targets: []target.Target{
				{Dir: ".", Artefact: "artefact-1"},
				{Dir: ".", Artefact: "artefact-2"},
			},
			Expected: []target.Target{
				{Dir: ".", Artefact: "artefact-1"},
				{Dir: ".", Artefact: "artefact-2"},
			},
		},
		{
			Name: "should add linked targets",
			Targets: []target.Target{
				{Dir: ".", Artefact: "artefact-1"},
				{Dir: ".", Artefact: "artefact-2"},
			},
			Links: []Link{
				{
					Source: target.Target{Dir: ".", Artefact: "artefact-1"},
					Target: target.Target{Dir: ".", Artefact: "artefact-3"},
				},
				{
					Source: target.Target{Dir: ".", Artefact: "artefact-2"},
					Target: target.Target{Dir: "subdir", Artefact: "artefact-1"},
				},
			},
			Expected: []target.Target{
				{Dir: ".", Artefact: "artefact-1"},
				{Dir: ".", Artefact: "artefact-2"},
				{Dir: ".", Artefact: "artefact-3"},
				{Dir: "subdir", Artefact: "artefact-1"},
			},
		},
		{
			Name: "should eliminate input and linked duplicates",
			Targets: []target.Target{
				{Dir: ".", Artefact: "artefact-1"},
				{Dir: ".", Artefact: "artefact-2"},
				{Dir: ".", Artefact: "artefact-2"},
				{Dir: ".", Artefact: "artefact-3"},
			},
			Links: []Link{
				{
					Source: target.Target{Dir: ".", Artefact: "artefact-1"},
					Target: target.Target{Dir: ".", Artefact: "artefact-2"},
				},
				{
					Source: target.Target{Dir: ".", Artefact: "artefact-1"},
					Target: target.Target{Dir: "subdir", Artefact: "artefact-1"},
				},
				{
					Source: target.Target{Dir: ".", Artefact: "artefact-3"},
					Target: target.Target{Dir: "subdir", Artefact: "artefact-1"},
				},
			},
			Expected: []target.Target{
				{Dir: ".", Artefact: "artefact-1"},
				{Dir: ".", Artefact: "artefact-2"},
				{Dir: ".", Artefact: "artefact-3"},
				{Dir: "subdir", Artefact: "artefact-1"},
			},
		},
	} {
		t.Run(test.Name, func(u *testing.T) {
			deduped := getDedupedTargets(test.Targets, test.Links)

			if test.Expected != nil {
				if deduped != nil {
					assert.Equal(u, test.Expected, deduped)
				} else {
					assert.Fail(u, "expected a list of targets")
				}
			}
		})
	}
}
