package storage

import (
	"encoding/json"
	"fmt"
	"github.com/bloxapp/dkg/dkg"
	"os"
	"path/filepath"
)

const (
	NodeFilename = "node"
)

func SaveNodeToDisk(node *dkg.Node, basePath string, password string) error {
	secrets := NewSecretsFromNode(node, password)
	if err := secrets.SaveToDisk(basePath); err != nil {
		return err
	}

	path := filepath.Join(basePath, fmt.Sprintf("%s_%d", NodeFilename, node.Index))
	return SaveJson(path, node)
}

func LoadNodes(basePath string) ([]*dkg.Node, error) {
	files, err := os.ReadDir(basePath)
	if err != nil {
		return nil, err
	}

	var ret []*dkg.Node
	for _, file := range files {
		byts, err := ReadJson(filepath.Join(basePath, file.Name()))
		if err != nil {
			return nil, err
		}

		node := &dkg.Node{}
		if err := json.Unmarshal(byts, &node); err != nil {
			continue // might not be a node file
		}
		ret = append(ret, node)
	}

	return ret, nil
}
