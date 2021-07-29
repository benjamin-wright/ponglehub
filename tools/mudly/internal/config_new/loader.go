package config_new

import (
	"errors"
	"fmt"
	"path"
	"regexp"

	"ponglehub.co.uk/tools/mudly/internal/target"
)

func getDirs(targets []target.Target) []string {
	dirs := []string{}

	for _, target := range targets {
		dupe := false
		for _, dir := range dirs {
			if target.Dir == dir {
				dupe = true
				break
			}
		}

		if !dupe {
			dirs = append(dirs, target.Dir)
		}
	}

	return dirs
}

func LoadConfigs(targets []target.Target) ([]Config, error) {
	configs := []Config{}

	for _, dir := range getDirs(targets) {
		cfg, err := getConfigData(dir)
		if err != nil {
			return nil, err
		}

		configs = append(configs, cfg)
	}

	return configs, nil
}

func getConfigData(filepath string) (Config, error) {
	cfg := Config{}

	cfg.Path = filepath

	r, err := openFile(path.Join(filepath, "Mudfile"))
	if err != nil {
		return cfg, err
	}

	r.prune()

	for r.nextLine() {
		switch r.getLineType() {
		case ARTEFACT_LINE:
			artefact, err := getArtefact(r)
			if err != nil {
				return cfg, err
			}

			cfg.Artefacts = append(cfg.Artefacts, artefact)
		case PIPELINE_LINE:
			pipeline, err := getPipeline(r)
			if err != nil {
				return cfg, err
			}

			cfg.Pipelines = append(cfg.Pipelines, pipeline)
		case ENV_LINE:
			name, value, err := getEnv(r)
			if err != nil {
				return cfg, err
			}

			if cfg.Env == nil {
				cfg.Env = map[string]string{}
			}

			cfg.Env[name] = value
		default:
			return cfg, fmt.Errorf("unknown line type: %s", r.line())
		}
	}

	return cfg, nil
}

var artefactRegex *regexp.Regexp = regexp.MustCompile(`^ARTEFACT (\S+)$`)

func getArtefact(r *reader) (Artefact, error) {
	artefact := Artefact{}
	firstLine := r.line()
	matches := artefactRegex.FindStringSubmatch(firstLine)

	if len(matches) != 2 {
		return artefact, fmt.Errorf("failed to parse artefact line \"%s\", wrong number of arguments", firstLine)
	}

	return Artefact{
		Name: matches[1],
	}, nil
}

func getPipeline(r *reader) (Pipeline, error) {
	return Pipeline{}, errors.New("not implemented")
}

var envRegex *regexp.Regexp = regexp.MustCompile(`^(?:\s*)ENV (\S+)\=(\S+)$`)

func getEnv(r *reader) (string, string, error) {
	matches := envRegex.FindStringSubmatch(r.line())

	if matches == nil {
		return "", "", fmt.Errorf("env unknown syntax error for line \"%s\"", r.line())
	}

	if len(matches) != 3 {
		return "", "", fmt.Errorf("env match count error for line \"%s\" (found %d, expecting 2)", r.line(), len(matches)-1)
	}

	return matches[1], matches[2], nil
}
