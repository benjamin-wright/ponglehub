package services

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/geppetto/types"
)

type io interface {
	ReadJSON(path string) (map[string]interface{}, error)
	WriteJSON(path string, data map[string]interface{}) error
}

type commander interface {
	Run(workDir string, command string) (string, error)
}

// NPM collects methods related to NPM repos
type NPM struct {
	io  io
	cmd commander
}

// NewNpmService create a new NPM repo instance, or error if the path doesn't contain a nodejs project
func NewNpmService() NPM {
	return NPM{io: &IO{}, cmd: &Commander{}}
}

// Install run an NPM install
func (r NPM) Install(repo types.Repo) error {
	output, err := r.cmd.Run(repo.Path, "npm install --strict-ssl=false")
	if err != nil {
		return fmt.Errorf("Error installing NPM module:\nError\n%+v\nOutput:\n%s", err, output)
	}

	return nil
}

// Lint run the NPM lint script
func (r NPM) Lint(repo types.Repo) error {
	output, err := r.cmd.Run(repo.Path, "npm run lint")
	if err != nil {
		return fmt.Errorf("Error linting NPM module:\nError\n%+v\nOutput:\n%s", err, output)
	}

	return nil
}

// Test run the NPM test
func (r NPM) Test(repo types.Repo) error {
	output, err := r.cmd.Run(repo.Path, "npm test")
	if err != nil {
		return fmt.Errorf("Error testing NPM module:\nError\n%+v\nOutput:\n%s", err, output)
	}

	return nil
}

// Publish push the repo up to its registry
func (r NPM) Publish(repo types.Repo) error {
	output, err := r.cmd.Run(repo.Path, "npm publish --strict-ssl=false")
	if err != nil {
		return fmt.Errorf("Error installing NPM module:\nError\n%+v\nOutput:\n%s", err, output)
	}

	return nil
}

// GetLatestSHA get the SHA of the most recently published version of the module
func (r NPM) GetLatestSHA(repo types.Repo) (string, error) {
	return r.cmd.Run(repo.Path, "npm view --strict-ssl=false --json | jq '.dist.shasum' -r")
}

// GetCurrentSHA get the SHA of the current version of the module
func (r NPM) GetCurrentSHA(repo types.Repo) (string, error) {
	return r.cmd.Run(repo.Path, "npm publish --dry-run --json | jq '.shasum' -r")
}

// GetRepo returns a repo object representing the node project at the designated file path
func (r NPM) GetRepo(path string) (types.Repo, error) {
	empty := types.Repo{}

	packageJSON := path + "/package.json"

	data, err := r.io.ReadJSON(packageJSON)
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
		RepoType:  types.Node,
		DependsOn: []string{},
	}, nil
}

// GetDependencyNames returns an array containg the names of all this project's dependencies
func (r NPM) GetDependencyNames(repo types.Repo) ([]string, error) {
	packageJSON := repo.Path + "/package.json"

	data, err := r.io.ReadJSON(packageJSON)
	if err != nil {
		return nil, err
	}

	names := []string{}

	if deps, ok := data["dependencies"]; ok {
		depsMap, ok := deps.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Dependencies found, but format was wrong: %+v", deps)
		}

		for key := range depsMap {
			names = append(names, key)
		}
	}

	if deps, ok := data["devDependencies"]; ok {
		depsMap, ok := deps.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Dev dependencies found, but format was wrong: %+v", deps)
		}

		for key := range depsMap {
			names = append(names, key)
		}
	}

	return names, nil
}

// SetVersion update the version number in package.json
func (r NPM) SetVersion(repo types.Repo, version string) error {
	path := repo.Path + "/package.json"

	result, err := r.io.ReadJSON(path)
	if err != nil {
		return err
	}

	version, ok := result["version"].(string)
	if !ok {
		return errors.New("package.json did not include a 'version' field")
	}
	logrus.Infof("%s version: %s -> 1.0.0", result["name"], version)

	result["version"] = "1.0.0"

	return r.io.WriteJSON(path, result)
}
