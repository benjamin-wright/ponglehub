package services

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/geppetto/types"
)

type io interface {
	ReadJSON(path string) (map[string]interface{}, error)
	WriteJSON(path string, data map[string]interface{}) error
}

// NPM collects methods related to NPM repos
type NPM struct {
	io  io
	cmd commander
}

// NewNpmService create a new NPM service instance
func NewNpmService() NPM {
	return NPM{io: &IO{}, cmd: &Commander{}}
}

// GetRepo returns a repo object representing the node project at the designated file path
func (n NPM) GetRepo(path string) (types.Repo, error) {
	empty := types.Repo{}

	packageJSON := path + "/package.json"

	data, err := n.io.ReadJSON(packageJSON)
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

// Install run an NPM install
func (n NPM) Install(repo types.Repo) error {
	output, err := n.cmd.Run(repo.Path, "npm install --strict-ssl=false")
	if err != nil {
		return fmt.Errorf("Error installing NPM module:\nError\n%+v\nOutput:\n%s", err, output)
	}

	return nil
}

// Lint run the NPM lint script
func (n NPM) Lint(repo types.Repo) error {
	output, err := n.cmd.Run(repo.Path, "npm run lint --silent")
	if err != nil {
		return fmt.Errorf("Error linting NPM module:\nError\n%+v\nOutput:\n%s", err, output)
	}

	return nil
}

// Test run the NPM test
func (n NPM) Test(repo types.Repo) error {
	output, err := n.cmd.Run(repo.Path, "npm test --silent")
	if err != nil {
		return fmt.Errorf("Error testing NPM module:\nError\n%+v\nOutput:\n%s", err, output)
	}

	return nil
}

// Publish push the repo up to its registry
func (n NPM) Publish(repo types.Repo) error {
	output, err := n.cmd.Run(repo.Path, "npm publish --strict-ssl=false")
	if err != nil {
		return fmt.Errorf("Error installing NPM module:\nError\n%+v\nOutput:\n%s", err, output)
	}

	return nil
}

// GetLatestSHA get the SHA of the most recently published version of the module
func (n NPM) GetLatestSHA(repo types.Repo) (string, error) {
	return n.cmd.Run(repo.Path, "npm view --strict-ssl=false --json | jq '.dist.shasum' -r")
}

// GetCurrentSHA get the SHA of the current version of the module
func (n NPM) GetCurrentSHA(repo types.Repo) (string, error) {
	return n.cmd.Run(repo.Path, "npm publish --dry-run --json | jq '.shasum' -r")
}

// GetDependencyNames returns an array containg the names of all this project's dependencies
func (n NPM) GetDependencyNames(repo types.Repo) ([]string, error) {
	packageJSON := repo.Path + "/package.json"

	data, err := n.io.ReadJSON(packageJSON)
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
func (n NPM) SetVersion(repo types.Repo, version string) error {
	path := repo.Path + "/package.json"

	result, err := n.io.ReadJSON(path)
	if err != nil {
		return err
	}

	current, ok := result["version"].(string)
	if !ok {
		return errors.New("package.json did not include a 'version' field")
	}

	if version == "" {
		components := strings.Split(current, ".")
		patch, err := strconv.Atoi(components[2])
		if err != nil {
			return fmt.Errorf("Failed to convert patch version '%s' in semver: %s", components[2], current)
		}

		version = fmt.Sprintf("%s.%s.%d", components[0], components[1], patch+1)
	}

	logrus.Infof("%s version: %s -> %s", result["name"], current, version)

	result["version"] = version

	err = n.io.WriteJSON(path, result)
	if err != nil {
		return err
	}

	output, err := n.cmd.Run(repo.Path, "npx prettier-package-json --write ./package.json")
	if err != nil {
		return fmt.Errorf("Error linting NPM module:\nError\n%+v\nOutput:\n%s", err, output)
	}

	return nil
}
