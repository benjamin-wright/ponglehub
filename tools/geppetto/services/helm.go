package services

import (
	"fmt"

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
	output, err := h.cmd.Run(repo.Path, "helm dep update")
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
