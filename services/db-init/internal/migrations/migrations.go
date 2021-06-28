package migrations

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type MigrationData struct {
	ID   int
	Path string
}

func Load(path string) ([]MigrationData, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	data := []MigrationData{}

	for _, file := range files {
		parts := strings.Split(file.Name(), "__")
		if len(parts) < 2 {
			return nil, fmt.Errorf("failed to parse migration file name \"%s\": missing \"__\" separator", file.Name())
		}

		if !strings.HasPrefix(parts[0], "V") {
			return nil, fmt.Errorf("failed to parse migration file name \"%s\": version not prepended with \"V\"", file.Name())
		}

		version, err := strconv.Atoi(strings.TrimPrefix(parts[0], "V"))
		if err != nil {
			return nil, fmt.Errorf("failed to parse migration file name \"%s\": failed to parse version number: %+v", file.Name(), err)
		}

		data = append(data, MigrationData{
			ID:   version,
			Path: file.Name(),
		})
	}

	return data, nil
}
