package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	path "path/filepath"
	"strings"
)

// RepoType indicates the type of data in a repo
type RepoType string

const (
	// Node repo is an NPM module
	Node RepoType = "Node"
	// Go repo is a GOLANG module
	Go RepoType = "Go"
)

// Repo represents a code repo
type Repo struct {
	// Name a unique name for the dependency
	Name string
	// Path the location of the code on disk
	Path string
	// The kind of code in the repo
	RepoType RepoType
	// The paths of other repos one which this one depends
	DependsOn []string
}

// Config represents the app configuration
type Config struct {
	Repos    []Repo
	BasePath string
}

// Includes returns true if the config includes the named repo
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

	cfg.BasePath, err = path.Abs(path.Dir(filepath))
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
