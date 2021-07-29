package config_new

import (
	"fmt"
	"strings"
)

type reader struct {
	lines []string
	index int
}

func openFile(filepath string) (*reader, error) {
	data, err := fsInstance.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %+v", err)
	}

	lines := strings.Split(string(data), "\n")
	index := -1

	return &reader{
		lines: lines,
		index: index,
	}, nil
}

func (r *reader) prune() {
	output := []string{}

	for _, line := range r.lines {
		trimmed := strings.TrimSpace(line)

		if trimmed == "" {
			continue
		}

		output = append(output, line)
	}

	r.lines = output
}

type lineType int

const (
	ARTEFACT_LINE lineType = iota
	PIPELINE_LINE
	ENV_LINE
	STEP_LINE
	UNKNOWN_LINE
	READER_ERROR
)

func (r *reader) getLineType() lineType {
	if r.index < 0 || r.index >= len(r.lines) {
		return READER_ERROR
	}

	trimmed := strings.TrimSpace(r.lines[r.index])

	if strings.HasPrefix(trimmed, "ENV") {
		return ENV_LINE
	}

	if strings.HasPrefix(trimmed, "ARTEFACT") {
		return ARTEFACT_LINE
	}

	return UNKNOWN_LINE
}
func (r *reader) indent() int       { return 0 }
func (r *reader) getBlockEnd() int  { return 0 }
func (r *reader) isNewEntity() bool { return false }
func (r *reader) nextLine() bool {
	r.index++
	return r.index < len(r.lines)
}
func (r *reader) line() string {
	if r.index < 0 || r.index >= len(r.lines) {
		return ""
	}

	return r.lines[r.index]
}
