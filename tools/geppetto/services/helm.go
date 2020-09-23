package services

import (
	"fmt"
	"strconv"
	"strings"

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
func (h *Helm) Install(repo types.Repo) error {
	output, err := h.cmd.Run(repo.Path, "rm -rf tmpcharts && helm dep update")
	if err != nil {
		return fmt.Errorf("Error installing Helm dependencies:\nError\n%+v\nOutput:\n%s", err, output)
	}

	return nil
}

// Lint run lint on a helm repo
func (h *Helm) Lint(repo types.Repo) error {
	cmd := "helm lint --values values.yaml"
	if h.io.FileExists(repo.Path + "/lint-values.yaml") {
		cmd = cmd + " --values lint-values.yaml"
	}

	output, err := h.cmd.Run(repo.Path, cmd)
	if err != nil {
		return fmt.Errorf("Error running lint on helm chart:\nError\n%+v\nOutput:\n%s", err, output)
	}

	return nil
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
func (h *Helm) Publish(repo types.Repo, chartRepo string) error {
	output, err := h.cmd.Run(repo.Path, "helm push . "+chartRepo+" --insecure")
	if err != nil {
		return fmt.Errorf("Error publishing helm chart:\nError\n%+v\nOutput:\n%s", err, output)
	}

	return nil
}
