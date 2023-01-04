package storage

import (
	"encoding/json"
	"os"
)

func SaveJson(path string, obj interface{}) error {
	byts, err := json.MarshalIndent(obj, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, byts, os.ModeDir)
}

func ReadJson(path string) ([]byte, error) {
	return os.ReadFile(path)
}
