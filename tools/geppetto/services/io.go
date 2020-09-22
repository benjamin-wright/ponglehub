package services

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// IO mockable functions for interacting with the file system
type IO struct{}

// ReadJSON reads arbitrary data from a json file
func (io *IO) ReadJSON(path string) (map[string]interface{}, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	byteData, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal([]byte(byteData), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ReadYAML reads arbitrary data from a yaml file
func (io *IO) ReadYAML(path string) (map[string]interface{}, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	byteData, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = yaml.Unmarshal([]byte(byteData), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// WriteJSON writes arbitrary data to a json file
func (io *IO) WriteJSON(path string, data map[string]interface{}) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}

	byteData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, byteData, 0644)
}

// WriteYAML writes arbitrary data to a yaml file
func (io *IO) WriteYAML(path string, data map[string]interface{}) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}

	byteData, err := yaml.Marshal(data)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, byteData, 0644)
}

// FileExists returns true if the file at the designated path exists
func (io *IO) FileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	if info.IsDir() {
		return false
	}

	return true
}

// Walk walks the file tree rooted at root, calling walkFn for each file or directory in the tree, including root. All errors that arise visiting files and directories are filtered by walkFn. The files are walked in lexical order, which makes the output deterministic but means that for very large directories Walk can be inefficient. Walk does not follow symbolic links.
func (io *IO) Walk(targetDir string, walkFunc filepath.WalkFunc) error {
	return filepath.Walk(targetDir, walkFunc)
}
