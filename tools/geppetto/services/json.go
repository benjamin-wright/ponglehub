package services

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func readJSON(path string) (map[string]interface{}, error) {
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

func writeJSON(path string, data map[string]interface{}) error {
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
