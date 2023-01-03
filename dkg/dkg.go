package dkg

import (
	"crypto/rsa"
	"github.com/drand/kyber/share/dkg"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
)

type Node struct {
	Index         uint32
	drand         *dkg.DistKeyGenerator
	ecies         *ECIES
	EncryptionPK  *rsa.PublicKey
	GenerateShare *bls.SecretKey
}

func NewNode(index uint32, s dkg.Suite, pk *rsa.PublicKey) *Node {
	return &Node{
		Index:        index,
		ecies:        NewRandomECIES(s),
		EncryptionPK: pk,
	}
}

func (n *Node) SetupDrandWithConfig(c *dkg.Config) error {
	c.Longterm = n.ecies.Priv
	drand, err := dkg.NewDistKeyHandler(c)
	if err != nil {
		errors.Wrap(err, "could not generate drand DKG")
	}
	n.drand = drand
	return nil
}

func (n *Node) EncryptShare() ([]byte, error) {
	return Encrypt(n.EncryptionPK, n.GenerateShare.Serialize())
}
