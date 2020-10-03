package services

import (
	"fmt"

	"ponglehub.co.uk/geppetto/types"
)

type goIo interface {
	ReadModfile(path string) (map[string]interface{}, error)
	FileExists(path string) bool
}

// Golang collects methods related to golang repos
type Golang struct {
	io  goIo
	cmd commander
}

// NewGolangService create a new Golang repo instance, or error if the path doesn't contain a golang project
func NewGolangService() Golang {
	return Golang{io: &IO{}, cmd: &Commander{}}
}

// GetRepo returns a repo object representing the node project at the designated file path
func (g Golang) GetRepo(path string) (types.Repo, error) {
	empty := types.Repo{}

	moduleFile := path + "/go.mod"

	data, err := g.io.ReadModfile(moduleFile)
	if err != nil {
		return empty, err
	}

	name, ok := data["name"]
	if !ok {
		return empty, fmt.Errorf("Failed to read name from package.json: %s", path)
	}

	nameString, ok := name.(string)
	if !ok {
		return empty, fmt.Errorf("Failed to read name from package.json: %v", name)
	}

	return types.Repo{
		Name:      nameString,
		Path:      path,
		RepoType:  types.Golang,
		DependsOn: []string{},
	}, nil
}

// GetDependencyNames returns an array containg the names of all this project's dependencies
func (g *Golang) GetDependencyNames(repo types.Repo) ([]string, error) {
	moduleFile := repo.Path + "/go.mod"
	data, err := g.io.ReadModfile(moduleFile)
	if err != nil {
		return nil, err
	}

	return data["dependencies"].([]string), nil
}

// Tidy tidies up module dependencies for a golang repo
func (g *Golang) Tidy(repo types.Repo) error {
	output, err := g.cmd.Run(repo.Path, "go mod tidy")
	if err != nil {
		return fmt.Errorf("Error tidying go mod:\nError\n%+v\nOutput:\n%s", err, output)
	}

	return nil
}

// Install caches module dependencies for a golang repo
func (g *Golang) Install(repo types.Repo) error {
	output, err := g.cmd.Run(repo.Path, "go mod download")
	if err != nil {
		return fmt.Errorf("Error downloading go mod dependencies:\nError\n%+v\nOutput:\n%s", err, output)
	}

	return nil
}

// Test runs unit tests for a golang repo
func (g *Golang) Test(repo types.Repo) error {
	output, err := g.cmd.Run(repo.Path, "go test ./...")
	if err != nil {
		return fmt.Errorf("Error running unit tests:\nError\n%+v\nOutput:\n%s", err, output)
	}

	return nil
}

// Buildable returns true if the go repo is an application, or false if it is a library
func (g *Golang) Buildable(repo types.Repo) bool {
	return g.io.FileExists(repo.Path + "/main.go")
}

// Build builds the binary for a golang repo
func (g *Golang) Build(repo types.Repo) error {
	output, err := g.cmd.Run(repo.Path, "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/"+repo.Name)
	if err != nil {
		return fmt.Errorf("Error building :\nError\n%+v\nOutput:\n%s", err, output)
	}

	return nil
}
