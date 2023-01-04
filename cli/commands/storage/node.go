package storage

import (
	"fmt"
	"github.com/bloxapp/dkg/dkg"
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
