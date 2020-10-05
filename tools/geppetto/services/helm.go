package services

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"ponglehub.co.uk/geppetto/types"
)

type helmIo interface {
	FileExists(path string) bool
	ReadYAML(path string) (map[string]interface{}, error)
	WriteYAML(path string, data map[string]interface{}) error
}

// Helm collects methods related to Helm repos
type Helm struct {
	io  helmIo
	cmd commander
}

// NewHelmService create a new Helm service instance
func NewHelmService() Helm {
	return Helm{
		cmd: &Commander{},
		io:  &IO{},
	}
}

// GetRepo returns a repo object representing the node project at the designated file path
func (h *Helm) GetRepo(path string) (types.Repo, error) {
	empty := types.Repo{}

	chartYAML := path + "/Chart.yaml"

	data, err := h.io.ReadYAML(chartYAML)
	if err != nil {
		return empty, err
	}

	name, ok := data["name"]
	if !ok {
		return empty, fmt.Errorf("Failed to read name from Chart.yaml: %s", path)
	}

	nameString, ok := name.(string)
	if !ok {
		return empty, fmt.Errorf("Failed to read name from Chart.yaml: %v", name)
	}

	return types.Repo{
		Name:      nameString,
		Path:      path,
		RepoType:  types.Helm,
		DependsOn: []string{},
	}, nil
}

// Install install dependencies for a helm repo
func (h *Helm) Install(ctx context.Context, repo types.Repo) error {
	output, err := h.cmd.Run(ctx, repo.Path, "rm -rf tmpcharts && helm repo update && helm dep update")
	if err != nil {
		return fmt.Errorf("Error installing Helm dependencies:\nError\n%+v\nOutput:\n%s", err, output)
	}

	return nil
}

// Lint run lint on a helm repo
func (h *Helm) Lint(ctx context.Context, repo types.Repo) error {
	cmd := "helm lint --values values.yaml"
	if h.io.FileExists(repo.Path + "/lint-values.yaml") {
		cmd = cmd + " --values lint-values.yaml"
	}

	output, err := h.cmd.Run(ctx, repo.Path, cmd)
	if err != nil {
		return fmt.Errorf("Error running lint on helm chart:\nError\n%+v\nOutput:\n%s", err, output)
	}

	return nil
}

// GetDependencyNames returns an array containg the names of all this project's dependencies
func (h *Helm) GetDependencyNames(repo types.Repo) ([]string, error) {
	chartYAML := repo.Path + "/Chart.yaml"

	data, err := h.io.ReadYAML(chartYAML)
	if err != nil {
		return nil, err
	}

	names := []string{}

	if deps, ok := data["dependencies"]; ok {
		if deps == nil {
			return []string{}, nil
		}

		depsList, ok := deps.([]interface{})
		if !ok {
			return nil, fmt.Errorf("Dependencies found, but format was wrong: %+v", deps)
		}

		logrus.Infof("Got %d deps", len(depsList))

		for _, dep := range depsList {
			depMap, ok := dep.(map[interface{}]interface{})
			if !ok {
				logrus.Infof("dependency is not a map[string]interface{}: %+v", dep)
				continue
			}

			depName, ok := depMap["name"].(string)
			if !ok {
				logrus.Infof("dependency name is not a string: %+v", dep)
				continue
			}

			repository, ok := depMap["repository"].(string)
			if !ok {
				logrus.Infof("dependency repository is not a string: %+v", dep)
				continue
			}

			logrus.Infof("Repo %s", repository)

			if repository != "@local" {
				logrus.Infof("repository for %s is not @local: %s", depName, repository)
				continue
			}

			names = append(names, depName)
		}
	}

	return names, nil
}

// GetCurrentVersion gets the local version of the chart
func (h *Helm) GetCurrentVersion(repo types.Repo) (string, error) {
	chartYAML := repo.Path + "/Chart.yaml"

	data, err := h.io.ReadYAML(chartYAML)
	if err != nil {
		return "", err
	}

	current := data["version"].(string)

	return current, nil
}

// GetLatestVersion gets the most up-to-date version of the published chart
func (h *Helm) GetLatestVersion(ctx context.Context, repo types.Repo, chartRepo string) (string, error) {
	output, err := h.cmd.Run(ctx, repo.Path, fmt.Sprintf("helm show chart %s/%s", chartRepo, repo.Name))
	if err != nil {
		return "", fmt.Errorf("Error fetching latest version of helm chart:\nError\n%+v\nOutput:\n%s", err, output)
	}

	var result map[string]interface{}
	err = yaml.Unmarshal([]byte(output), &result)
	if err != nil {
		return "", err
	}

	return result["version"].(string), nil
}

// SetVersion bump the version of the chart in the local registry
func (h *Helm) SetVersion(repo types.Repo, version string) error {
	chartYAML := repo.Path + "/Chart.yaml"

	data, err := h.io.ReadYAML(chartYAML)
	if err != nil {
		return err
	}

	current := data["version"].(string)

	if version == "" {
		components := strings.Split(current, ".")
		patch, err := strconv.Atoi(components[2])
		if err != nil {
			return fmt.Errorf("Failed to convert patch version '%s' in semver: %s", components[2], current)
		}

		version = fmt.Sprintf("%s.%s.%d", components[0], components[1], patch+1)
	}

	data["version"] = version

	return h.io.WriteYAML(chartYAML, data)
}

// Publish publish the chart to a local registry
func (h *Helm) Publish(ctx context.Context, repo types.Repo, chartRepo string) error {
	output, err := h.cmd.Run(ctx, repo.Path, "helm push . "+chartRepo)
	if err != nil {
		return fmt.Errorf("Error publishing helm chart:\nError\n%+v\nOutput:\n%s", err, output)
	}

	return nil
}
