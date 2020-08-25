package config

import (
	"fmt"
)

// FileStruct struct for umarshalling config data
type FileStruct struct {
	Node []RepoStruct `json:"node"`
	Go   []RepoStruct `json:"go"`
}

// RepoStruct struct for unmarshalling config data
type RepoStruct struct {
	Name         string   `json:"name"`
	Path         string   `json:"path"`
	Dependencies []string `json:"dependencies"`
}

func (f FileStruct) toConfig() (Config, error) {
	cfg := Config{
		Repos: []Repo{},
	}

	for _, repo := range f.Node {
		if cfg.Includes(repo.Name) {
			return Config{}, fmt.Errorf("Failed to load config: repo with duplicate name (%s)", repo.Name)
		}

		cfg.Repos = append(cfg.Repos, Repo{
			Name:      repo.Name,
			Path:      repo.Path,
			RepoType:  Node,
			DependsOn: repo.Dependencies,
		})
	}

	for _, repo := range f.Go {
		if cfg.Includes(repo.Name) {
			return Config{}, fmt.Errorf("Failed to load config: repo with duplicate name (%s)", repo.Name)
		}

		cfg.Repos = append(cfg.Repos, Repo{
			Name:      repo.Name,
			Path:      repo.Path,
			RepoType:  Go,
			DependsOn: repo.Dependencies,
		})
	}

	return cfg, nil
}
