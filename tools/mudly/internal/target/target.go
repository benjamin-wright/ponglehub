package target

import (
	"errors"
	"fmt"
	"strings"
)

type Target struct {
	Dir      string
	Artefact string
}

func ParseTarget(path string) (*Target, error) {
	if path == "" {
		return nil, errors.New("failed to parse target with empty string")
	}

	parts := strings.Split(path, "+")
	if len(parts) != 2 {
		return nil, fmt.Errorf("failed to parse target from path: %s", path)
	}

	dir := parts[0]
	artefact := parts[1]

	if dir == "" {
		dir = "."
	}

	return &Target{
		Dir:      dir,
		Artefact: artefact,
	}, nil
}