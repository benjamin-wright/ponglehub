package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

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

func (c *Config) HasCircularDependencies() bool {
	built := []string{}
	building := true

	isBuilt := func(name string) bool {
		for _, b := range built {
			if name == b {
				return true
			}
		}

		return false
	}

	for building {
		building = false
		builds := []string{}

		for _, r := range c.Repos {
			if isBuilt(r.Name) {
				continue
			}

			depsBuilt := true
			for _, d := range r.DependsOn {
				if !isBuilt(d) {
					depsBuilt = false
				}
			}

			if depsBuilt {
				building = true
				builds = append(builds, r.Name)
			}
		}
	}

	if len(built) != len(c.Repos) {
		return true
	}

	return false
}

func (c *Config) Includes(name string) bool {
	if c == nil {
		return false
	}

	for _, repo := range c.Repos {
		if repo.Name == name {
			return true
		}
	}

	return false
}

func (c *Config) String() string {
	if c == nil {
		return "Nil config"
	}

	if len(c.Repos) == 0 {
		return "Empty"
	}

	var sb strings.Builder
	var first = true

	for _, repo := range c.Repos {
		if !first {
			sb.WriteString(", ")
		}

		sb.WriteString(fmt.Sprintf("[%s: %s]", repo.Name, repo.RepoType))

		if first {
			first = false
		}
	}

	return sb.String()
}

// FromFile create a new config object from the config file
func FromFile(filepath string) (*Config, error) {
	jsonFile, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var fileData FileStruct
	err = json.Unmarshal(byteValue, &fileData)
	if err != nil {
		return nil, err
	}

	cfg, err := fileData.toConfig()
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
