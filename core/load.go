package core

import (
	"os"
)

func yamlFromFile(filepath string) ([]byte, error) {
	readFile, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	return readFile, nil
}
