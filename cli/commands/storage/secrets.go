package storage

import (
	"encoding/json"
	"github.com/bloxapp/dkg/dkg"
	"github.com/pkg/errors"
	"path/filepath"
)

const (
	SecretsFilename = "secrets"

	SessionSKKey        = "SessionSKKey"
	GeneratedShareSKKey = "GeneratedShareSK"
)

type Secrets map[string]string

func NewSecretsFromNode(n *dkg.Node, password string) Secrets {
	ret := Secrets{}

	if n.Ecies.GetPrivateKey() != nil {
		byts, err := n.Ecies.GetPrivateKey().MarshalBinary()
		if err != nil {
			panic(err.Error())
		}
		if err := ret.SetWithPassword(byts, SessionSKKey, password); err != nil {
			panic(err.Error())
		}
	}

	if n.DidSetGeneratedShare() {
		byts := n.GetGenerateShare().Serialize()
		if err := ret.SetWithPassword(byts, GeneratedShareSKKey, password); err != nil {
			panic(err.Error())
		}
	}

	return ret
}

func (s Secrets) SaveToDisk(basePath string) error {
	path := filepath.Join(basePath, SecretsFilename)
	return SaveJson(path, s)
}

func (s Secrets) LoadFromDisk(basePath string) error {
	path := filepath.Join(basePath, SecretsFilename)
	byts, err := ReadJson(path)
	if err != nil {
		return err
	}

	return json.Unmarshal(byts, &s)
}

func (s Secrets) SetWithPassword(data []byte, key, password string) error {
	ciper, err := encrypt(password, data)
	if err != nil {
		return err
	}

	s[key] = string(ciper)
	return nil
}

func (s Secrets) ReadWithPassword(key, password string) ([]byte, error) {
	val, ok := s[key]
	if !ok {
		return nil, errors.New("key not found")
	}
	return decrypt(password, []byte(val))
}
