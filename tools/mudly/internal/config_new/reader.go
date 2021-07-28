package config_new

import (
	"fmt"
	"strings"
)

type dockerfile struct {
	name    string
	content string
}

type configData struct {
	dockerfile []dockerfile
	artefacts  []artefactData
	pipelines  []pipelineData
	env        map[string]string
}

type pipelineData struct {
	name  string
	steps []stepData
}

type artefactData struct {
	name      string
	dependsOn []string
	env       map[string]string
	steps     []stepData
	pipeline  string
}

type stepData struct {
	name       string
	env        map[string]string
	watch      []string
	command    string
	dockerfile string
}

func getConfigData(filepath string) (configData, error) {
	cfg := configData{}

	data, err := fsInstance.ReadFile(filepath)
	if err != nil {
		return cfg, fmt.Errorf("failed to open config file: %+v", err)
	}

	lines := strings.Split(string(data), "\n")
	index := 0

	return cfg, nil
}
