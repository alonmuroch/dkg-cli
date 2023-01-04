package dkg

import (
	"fmt"
	"github.com/drand/kyber"
	bls3 "github.com/drand/kyber-bls12381"
	"github.com/drand/kyber/share"
	"github.com/drand/kyber/share/dkg"
	bls2 "github.com/drand/kyber/sign/bls"
	"github.com/drand/kyber/util/random"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/pkg/errors"
	"math/rand"
)

func InitBLS() {
	_ = bls.Init(bls.BLS12_381)
	_ = bls.SetETHmode(bls.EthModeDraft07)
}

var Suite = bls3.NewBLS12381Suite()
var AuthScheme = bls2.NewSchemeOnG1(Suite)

type ECIES struct {
	Pub   []byte
	priv  []byte
	suite dkg.Suite
}

func NewRandomECIES(s dkg.Suite) *ECIES {
	private := s.Scalar().Pick(random.New())
	pkByts, err := private.MarshalBinary()
	if err != nil {
		panic(err.Error())
	}
	pubByts, err := s.Point().Mul(private, nil).MarshalBinary()
	if err != nil {
		panic(err.Error())
	}

	return &ECIES{
		suite: s,
		priv:  pkByts,
		Pub:   pubByts,
	}
}

func (e *ECIES) GetPublicKey() kyber.Point {
	ret := e.suite.Point()
	if err := ret.UnmarshalBinary(e.Pub); err != nil {
		panic(err.Error())
	}
	return ret
}

func (e *ECIES) GetPrivateKey() kyber.Scalar {
	ret := e.suite.Scalar()
	if err := ret.UnmarshalBinary(e.priv); err != nil {
		panic(err.Error())
	}
	return ret
}

// NonceLength is the length of the nonce
const NonceLength = 32

// GetNonce returns a suitable nonce to feed in the DKG config.
func GetNonce() []byte {
	var nonce [NonceLength]byte
	n, err := rand.Read(nonce[:])
	if n != NonceLength {
		panic("could not read enough random bytes for nonce")
	}
	if err != nil {
		panic(err)
	}
	return nonce[:]
}

func resultToShareSecretKey(result *dkg.Result) (*bls.SecretKey, error) {
	share := result.Key.PriShare()
	bytsSk, err := share.V.MarshalBinary()
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal share")
	}
	sk := &bls.SecretKey{}
	if err := sk.Deserialize(bytsSk); err != nil {
		return nil, errors.Wrap(err, "could not deserialized secret key")
	}
	return sk, nil
}

func resultsToValidatorPK(results []*dkg.Result, suite dkg.Suite) (*bls.PublicKey, error) {
	exp := share.NewPubPoly(suite, suite.Point().Base(), results[0].Key.Commitments())
	bytsPK, err := exp.Eval(0).V.MarshalBinary()
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal share")
	}
	pk := &bls.PublicKey{}
	if err := pk.Deserialize(bytsPK); err != nil {
		return nil, errors.Wrap(err, "could not deserialized public key")
	}
	return pk, nil
}

// ReconstructSignatures receives a map of user indexes and serialized bls.Sign.
// It then reconstructs the original threshold signature using lagrange interpolation
func reconstructSignatures(signatures map[int][]byte) (*bls.Sign, error) {
	reconstructedSig := bls.Sign{}

	idVec := make([]bls.ID, 0)
	sigVec := make([]bls.Sign, 0)

	for index, signature := range signatures {
		blsID := bls.ID{}
		err := blsID.SetDecString(fmt.Sprintf("%d", index))
		if err != nil {
			return nil, err
		}

		idVec = append(idVec, blsID)
		blsSig := bls.Sign{}

		err = blsSig.Deserialize(signature)
		if err != nil {
			return nil, err
		}

		sigVec = append(sigVec, blsSig)
	}
	err := reconstructedSig.Recover(sigVec, idVec)
	return &reconstructedSig, err
}
