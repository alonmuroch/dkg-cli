package dkg

import (
	"crypto/rsa"
	"github.com/drand/kyber/share/dkg"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
)

type Node struct {
	Index        uint32
	Ecies        *ECIES
	EncryptionPK []byte

	generatedShare []byte

	drand *dkg.DistKeyGenerator
}

func NewNode(index uint32, s dkg.Suite, pk []byte) *Node {
	return &Node{
		Index:        index,
		Ecies:        NewRandomECIES(s),
		EncryptionPK: pk,
	}
}

func (n *Node) SetupDrandWithConfig(c *dkg.Config) error {
	c.Longterm = n.Ecies.GetPrivateKey()
	drand, err := dkg.NewDistKeyHandler(c)
	if err != nil {
		errors.Wrap(err, "could not generate drand DKG")
	}
	n.drand = drand
	return nil
}

func (n *Node) getEncryptionPK() *rsa.PublicKey {
	encryptionPK, err := PemToPublicKey(n.EncryptionPK)
	if err != nil {
		panic(err.Error())
	}
	return encryptionPK
}

func (n *Node) DidSetGeneratedShare() bool {
	return len(n.generatedShare) > 0
}

func (n *Node) GetGenerateShare() *bls.SecretKey {
	ret := &bls.SecretKey{}
	if err := ret.Deserialize(n.generatedShare); err != nil {
		panic(err.Error())
	}
	return ret
}

func (n *Node) EncryptShare() ([]byte, error) {
	return Encrypt(n.getEncryptionPK(), n.generatedShare)
}
