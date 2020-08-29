package services

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/geppetto/types"
)

type fileIO interface {
	readJSON(path string) (map[string]interface{}, error)
	writeJSON(path string, data map[string]interface{}) error
}

type commander interface {
	run(workDir string, command string) (string, error)
}

// NPM collects methods related to NPM repos
type NPM struct {
	io  fileIO
	cmd commander
}

// NewNpmService create a new NPM repo instance, or error if the path doesn't contain a nodejs project
func NewNpmService() NPM {
	return NPM{io: &defaultIO{}, cmd: &defaultCommander{}}

	// packageJSON := path + "/package.json"

	// info, err := os.Stat(packageJSON)
	// if err != nil {
	// 	return empty, err
	// }

	// if info.IsDir() {
	// 	return empty, fmt.Errorf("Expected %s to be a file, not a directory", packageJSON)
	// }
}

// Install run an NPM install
func (r NPM) Install(repo types.Repo) error {
	output, err := r.cmd.run(repo.Path, "npm install --strict-ssl=false")
	if err != nil {
		return fmt.Errorf("Error installing NPM module:\nError\n%+v\nOutput:\n%s", err, output)
	}

	return nil
}

// Lint run the NPM lint script
func (r NPM) Lint(repo types.Repo) error {
	output, err := r.cmd.run(repo.Path, "npm run lint")
	if err != nil {
		return fmt.Errorf("Error linting NPM module:\nError\n%+v\nOutput:\n%s", err, output)
	}

	return nil
}

// Test run the NPM test
func (r NPM) Test(repo types.Repo) error {
	output, err := r.cmd.run(repo.Path, "npm test")
	if err != nil {
		return fmt.Errorf("Error testing NPM module:\nError\n%+v\nOutput:\n%s", err, output)
	}

	return nil
}

// Publish push the repo up to its registry
func (r NPM) Publish(repo types.Repo) error {
	output, err := r.cmd.run(repo.Path, "npm publish --strict-ssl=false")
	if err != nil {
		return fmt.Errorf("Error installing NPM module:\nError\n%+v\nOutput:\n%s", err, output)
	}

	return nil
}

// GetLatestSHA get the SHA of the most recently published version of the module
func (r NPM) GetLatestSHA(repo types.Repo) (string, error) {
	return r.cmd.run(repo.Path, "npm view --strict-ssl=false --json | jq '.dist.shasum' -r")
}

// GetCurrentSHA get the SHA of the current version of the module
func (r NPM) GetCurrentSHA(repo types.Repo) (string, error) {
	return r.cmd.run(repo.Path, "npm publish --dry-run --json | jq '.shasum' -r")
}

// SetVersion update the version number in package.json
func (r NPM) SetVersion(repo types.Repo, version string) error {
	path := repo.Path + "/package.json"

	result, err := r.io.readJSON(path)
	if err != nil {
		return err
	}

	version, ok := result["version"].(string)
	if !ok {
		return errors.New("package.json did not include a 'version' field")
	}
	logrus.Infof("%s version: %s -> 1.0.0", result["name"], version)

	result["version"] = "1.0.0"

	return r.io.writeJSON(path, result)
}
