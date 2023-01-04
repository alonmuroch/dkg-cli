package storage

import (
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

	path := filepath.Join(basePath, NodeFilename)
	return SaveJson(path, node)
}
