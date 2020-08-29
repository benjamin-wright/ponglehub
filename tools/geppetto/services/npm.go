package services

import (
	"errors"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

type fileIO interface {
	readJSON(path string) (map[string]interface{}, error)
	writeJSON(path string, data map[string]interface{}) error
}

type commander interface {
	run(workDir string, command string) (string, error)
}

// NPMRepo collects methods related to NPM repos
type NPMRepo struct {
	path string
	io   fileIO
	cmd  commander
}

// NewNpmRepo create a new NPM repo instance, or error if the path doesn't contain a nodejs project
func NewNpmRepo(path string) (NPMRepo, error) {
	repo := NPMRepo{path: path, io: &defaultIO{}, cmd: &defaultCommander{}}
	empty := NPMRepo{}

	packageJSON := path + "/package.json"

	info, err := os.Stat(packageJSON)
	if err != nil {
		return empty, err
	}

	if info.IsDir() {
		return empty, fmt.Errorf("Expected %s to be a file, not a directory", packageJSON)
	}

	return repo, nil
}

// Install run an NPM install
func (r NPMRepo) Install() error {
	output, err := r.cmd.run(r.path, "npm install --strict-ssl=false")
	if err != nil {
		return fmt.Errorf("Error installing NPM module:\nError\n%+v\nOutput:\n%s", err, output)
	}

	return nil
}

// Lint run the NPM lint script
func (r NPMRepo) Lint() error {
	output, err := r.cmd.run(r.path, "npm run lint")
	if err != nil {
		return fmt.Errorf("Error linting NPM module:\nError\n%+v\nOutput:\n%s", err, output)
	}

	return nil
}

// Test run the NPM test
func (r NPMRepo) Test() error {
	output, err := r.cmd.run(r.path, "npm test")
	if err != nil {
		return fmt.Errorf("Error testing NPM module:\nError\n%+v\nOutput:\n%s", err, output)
	}

	return nil
}

// Publish push the repo up to its registry
func (r NPMRepo) Publish() error {
	output, err := r.cmd.run(r.path, "npm publish --strict-ssl=false")
	if err != nil {
		return fmt.Errorf("Error installing NPM module:\nError\n%+v\nOutput:\n%s", err, output)
	}

	return nil
}

// GetLatestSHA get the SHA of the most recently published version of the module
func (r NPMRepo) GetLatestSHA() (string, error) {
	return r.cmd.run(r.path, "npm view --strict-ssl=false --json | jq '.dist.shasum' -r")
}

// GetCurrentSHA get the SHA of the current version of the module
func (r NPMRepo) GetCurrentSHA() (string, error) {
	return r.cmd.run(r.path, "npm publish --dry-run --json | jq '.shasum' -r")
}

// SetVersion update the version number in package.json
func (r NPMRepo) SetVersion(version string) error {
	path := r.path + "/package.json"

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
