package services

import (
	"fmt"
	"os"
)

// NPMRepo collects methods related to NPM repos
type NPMRepo struct {
	path string
}

// NewNpmRepo create a new NPM repo instance, or error if the path doesn't contain a nodejs project
func NewNpmRepo(path string) (NPMRepo, error) {
	repo := NPMRepo{path: path}
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
	output, err := run(r.path, "npm install --strict-ssl=false")
	if err != nil {
		return fmt.Errorf("Error installing NPM module:\nError\n%+v\nOutput:\n%s", err, output)
	}

	return nil
}

// Lint run the NPM lint script
func (r NPMRepo) Lint() error {
	output, err := run(r.path, "npm run lint")
	if err != nil {
		return fmt.Errorf("Error linting NPM module:\nError\n%+v\nOutput:\n%s", err, output)
	}

	return nil
}

// Test run the NPM test
func (r NPMRepo) Test() error {
	output, err := run(r.path, "npm test")
	if err != nil {
		return fmt.Errorf("Error testing NPM module:\nError\n%+v\nOutput:\n%s", err, output)
	}

	return nil
}

// Publish push the repo up to its registry
func (r NPMRepo) Publish() error {
	output, err := run(r.path, "npm publish --strict-ssl=false")
	if err != nil {
		return fmt.Errorf("Error installing NPM module:\nError\n%+v\nOutput:\n%s", err, output)
	}

	return nil
}
