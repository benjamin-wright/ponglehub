package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// FileStruct struct for umarshalling config data
type FileStruct struct {
	Node []RepoStruct `json:"node"`
}

// RepoStruct struct for unmarshalling config data
type RepoStruct struct {
	Name         string   `json:"name"`
	Path         string   `json:"path"`
	Dependencies []string `json:"dependencies"`
}

func (f FileStruct) toConfig() Config {
	cfg := Config{
		Repos: []Repo{},
	}

	for _, repo := range f.Node {
		cfg.Repos = append(cfg.Repos, Repo{
			Name:      repo.Name,
			Path:      repo.Path,
			RepoType:  Node,
			DependsOn: repo.Dependencies,
		})
	}

	return cfg
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

	cfg := fileData.toConfig()

	return &cfg, nil
}
