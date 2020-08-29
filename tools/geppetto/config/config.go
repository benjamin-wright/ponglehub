package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// Config represents the app configuration
type Config struct {
	Ingore []string `json:"ignore"`
}

// FromFile create a new config object from the config file
func FromFile(filepath string) (*Config, error) {
	jsonFile, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = json.Unmarshal(byteValue, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
