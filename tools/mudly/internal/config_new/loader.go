package config_new

import (
	"errors"
	"fmt"
	"path"
	"regexp"
	"strings"

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

func getArtefact(r *reader) (Artefact, error) {
	artefact := Artefact{}
	firstLine := r.line()

	trimmed := strings.TrimSpace(firstLine)
	parts := strings.Split(trimmed, " ")

	if len(parts) != 2 {
		return artefact, fmt.Errorf("failed to parse artefact line \"%s\", wrong number of arguments", firstLine)
	}

	artefact.Name = parts[1]

	targetIndent := r.indent()

	for r.nextLine() {
		indent := r.indent()

		if indent <= targetIndent {
			r.previousLine()
			break
		}

		switch r.getLineType() {
		case ENV_LINE:
			name, value, err := getEnv(r)
			if err != nil {
				return artefact, err
			}

			if artefact.Env == nil {
				artefact.Env = map[string]string{}
			}

			artefact.Env[name] = value
		case PIPELINE_LINE:
			name, err := getPipelineLink(r)
			if err != nil {
				return artefact, err
			}

			artefact.Pipeline = name
		case DEPENDS_LINE:
			t, err := getDepends(r)
			if err != nil {
				return artefact, err
			}

			if artefact.DependsOn == nil {
				artefact.DependsOn = []target.Target{}
			}

			artefact.DependsOn = append(artefact.DependsOn, t)
		case STEP_LINE:
			step, err := getStep(r)
			if err != nil {
				return artefact, err
			}

			if artefact.Steps == nil {
				artefact.Steps = []Step{}
			}

			artefact.Steps = append(artefact.Steps, step)
		default:
			return artefact, fmt.Errorf("unknown line type: %s", r.line())
		}
	}

	return artefact, nil
}

func getPipeline(r *reader) (Pipeline, error) {
	pipeline := Pipeline{}
	firstLine := r.line()

	trimmed := strings.TrimSpace(firstLine)
	parts := strings.Split(trimmed, " ")

	if len(parts) != 2 {
		return pipeline, fmt.Errorf("failed to parse pipeline line \"%s\", wrong number of arguments", firstLine)
	}

	pipeline.Name = parts[1]

	targetIndent := r.indent()

	for r.nextLine() {
		indent := r.indent()

		if indent <= targetIndent {
			r.previousLine()
			break
		}

		switch r.getLineType() {
		case ENV_LINE:
			name, value, err := getEnv(r)
			if err != nil {
				return pipeline, err
			}

			if pipeline.Env == nil {
				pipeline.Env = map[string]string{}
			}

			pipeline.Env[name] = value
		case STEP_LINE:
			step, err := getStep(r)
			if err != nil {
				return pipeline, err
			}

			if pipeline.Steps == nil {
				pipeline.Steps = []Step{}
			}

			pipeline.Steps = append(pipeline.Steps, step)
		default:
			return pipeline, fmt.Errorf("unknown line type: %s", r.line())
		}
	}

	return pipeline, nil
}

func getPipelineLink(r *reader) (string, error) {
	trimmed := strings.TrimSpace(r.line())

	parts := strings.Split(trimmed, " ")

	if len(parts) != 2 {
		return "", fmt.Errorf("pipeline unknown syntax error for line \"%s\"", r.line())
	}

	return parts[1], nil
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

func getDepends(r *reader) (target.Target, error) {
	trimmed := strings.TrimSpace(r.line())

	parts := strings.Split(trimmed, " ")

	if len(parts) != 3 {
		return target.Target{}, fmt.Errorf("depends unknown syntax error for line \"%s\"", r.line())
	}

	t, err := target.ParseTarget(parts[2])
	if err != nil {
		return target.Target{}, err
	}

	if t == nil {
		return target.Target{}, fmt.Errorf("expected a target but got nil: \"%s\"", r.line())
	}

	return *t, nil
}

func getStep(r *reader) (Step, error) {
	step := Step{}
	firstLine := r.line()

	trimmed := strings.TrimSpace(firstLine)
	parts := strings.Split(trimmed, " ")

	if len(parts) != 2 {
		return step, fmt.Errorf("failed to parse artefact line \"%s\", wrong number of arguments", firstLine)
	}

	step.Name = parts[1]

	targetIndent := r.indent()

	for r.nextLine() {
		indent := r.indent()

		if indent <= targetIndent {
			r.previousLine()
			break
		}

		switch r.getLineType() {
		case ENV_LINE:
			name, value, err := getEnv(r)
			if err != nil {
				return step, err
			}

			if step.Env == nil {
				step.Env = map[string]string{}
			}

			step.Env[name] = value
		case WATCH_LINE:
			paths, err := getWatchPaths(r)
			if err != nil {
				return step, err
			}

			if step.Watch == nil {
				step.Watch = []string{}
			}

			step.Watch = append(step.Watch, paths...)
		case CONDITION_LINE:
			condition, err := getStringOrMultiline(r)
			if err != nil {
				return step, err
			}

			step.Condition = condition
		case COMMAND_LINE:
			command, err := getStringOrMultiline(r)
			if err != nil {
				return step, err
			}

			step.Command = command
		case DOCKER_LINE:
			dockerfile, err := getStringOrMultiline(r)
			if err != nil {
				return step, err
			}

			step.Dockerfile = dockerfile
		default:
			return step, fmt.Errorf("unknown line type: %s", r.line())
		}
	}

	return step, nil
}

func getWatchPaths(r *reader) ([]string, error) {
	trimmed := strings.TrimSpace(r.line())

	parts := strings.Split(trimmed, " ")

	if len(parts) < 2 {
		return nil, fmt.Errorf("unknown syntax error for line \"%s\"", r.line())
	}

	return parts[1:], nil
}

func getStringOrMultiline(r *reader) (string, error) {
	trimmed := strings.TrimSpace(r.line())

	parts := strings.Split(trimmed, " ")

	if len(parts) > 1 {
		return strings.Join(parts[1:], " "), nil
	}

	lines := []string{}
	targetIndent := r.indent()
	steppedIndent := targetIndent

	for r.nextLine() {
		indent := r.indent()

		if indent <= targetIndent {
			r.previousLine()
			break
		}

		if steppedIndent == targetIndent {
			steppedIndent = indent
		}

		lines = append(lines, r.line()[steppedIndent:])
	}

	if len(lines) == 0 {
		return "", errors.New("empty string / multiline-string not supported")
	}

	return strings.Join(lines, "\n"), nil
}
