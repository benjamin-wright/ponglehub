package commands

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/geppetto/config"
)

// NpmCommand represents a command to NPM
type NpmCommand struct {
	basePath string
	repo     config.Repo
}

// MakeNpmCommand make a new NPM rollback command object
func MakeNpmCommand(basePath string, repo config.Repo) *NpmCommand {
	cmd := NpmCommand{
		basePath: basePath,
		repo:     repo,
	}

	return &cmd
}

// Run run the NPM command and return an error if it fails
func (e NpmCommand) Run() error {
	logrus.Debugf("Running rollback on %s", e.repo.Name)
	path := e.basePath + "/" + e.repo.Path + "/package.json"

	byteData, err := e.readFile(path)
	if err != nil {
		return err
	}

	var result map[string]interface{}
	err = json.Unmarshal([]byte(byteData), &result)
	if err != nil {
		return err
	}

	version, ok := result["version"]
	if !ok {
		return errors.New("package.json did not include a 'version' field")
	}
	logrus.Infof("%s version: %s -> 1.0.0", e.repo.Name, version)

	result["version"] = "1.0.0"

	return e.writeFile(path, result)
}

// Name the name of the job
func (e NpmCommand) Name() string {
	return e.repo.Name
}

func (e NpmCommand) readFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	byteData, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return byteData, nil
}

func (e NpmCommand) writeFile(path string, data map[string]interface{}) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}

	byteData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, byteData, 0644)
}
