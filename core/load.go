package core

import (
	"gopkg.in/yaml.v3"
	"os"
)

func LoadClashConfigFromFile(filePath string) (map[string]any, error) {
	var data map[string]any
	yamlData, err := yamlFromFile(filePath)
	if err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(yamlData, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func yamlFromFile(filepath string) ([]byte, error) {
	readFile, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	return readFile, nil
}
