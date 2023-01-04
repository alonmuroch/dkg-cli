package dkg

import (
	"encoding/hex"
	"fmt"
	"github.com/drand/kyber/share/dkg"
	bls3 "github.com/herumi/bls-eth-go-binary/bls"
	"github.com/stretchr/testify/require"
	"testing"
)

var TestSuite = Suite
var TestAuthScheme = AuthScheme

const (
	N = 4
	T = 3
)

var TestNodesEncryptionKeys = func() [][]byte {
	_, pk1Byts, _ := GenerateKey()
	_, pk2Byts, _ := GenerateKey()
	_, pk3Byts, _ := GenerateKey()
	_, pk4Byts, _ := GenerateKey()

	return [][]byte{pk1Byts, pk2Byts, pk3Byts, pk4Byts}
}()

var TestNodes = func() []*Node {
	return []*Node{
		NewNode(1, TestSuite.G1().(dkg.Suite), TestNodesEncryptionKeys[0]),
		NewNode(2, TestSuite.G1().(dkg.Suite), TestNodesEncryptionKeys[1]),
		NewNode(3, TestSuite.G1().(dkg.Suite), TestNodesEncryptionKeys[2]),
		NewNode(4, TestSuite.G1().(dkg.Suite), TestNodesEncryptionKeys[3]),
	}
}()

var TestWithdrawalCredentials = func() []byte {
	s := "005b55a6c968852666b132a80f53712e5097b0fca86301a16992e695a8e86f16"
	ret, _ := hex.DecodeString(s)
	return ret
}()
var TestDrandNodes = func() []dkg.Node {
	return []dkg.Node{
		{
			Index:  1,
			Public: TestNodes[0].Ecies.GetPublicKey(),
		},
		{
			Index:  2,
			Public: TestNodes[1].Ecies.GetPublicKey(),
		},
		{
			Index:  3,
			Public: TestNodes[2].Ecies.GetPublicKey(),
		},
		{
			Index:  4,
			Public: TestNodes[3].Ecies.GetPublicKey(),
		},
	}
}()

func reconstructSK(t *testing.T, sks []bls3.SecretKey) *bls3.SecretKey {
	s := &bls3.SecretKey{}
	require.NoError(t, s.Recover(
		sks,
		[]bls3.ID{blsID(t, 1), blsID(t, 2), blsID(t, 3), blsID(t, 4)},
	))
	return s
}

func TestDKGFull(t *testing.T) {
	InitBLS()

	nonce := GetNonce()
	for _, n := range TestNodes {
		require.NoError(t, n.SetupDrandWithConfig(&dkg.Config{
			Suite:     TestSuite.G1().(dkg.Suite),
			NewNodes:  TestDrandNodes,
			Threshold: T,
			Auth:      TestAuthScheme,
			Nonce:     nonce,
		}))
	}

	// Step 1
	var deals []*dkg.DealBundle
	for _, n := range TestNodes {
		d, err := n.drand.Deals()
		require.NoError(t, err)
		deals = append(deals, d)
	}
	require.NotEmpty(t, deals)

	// Step 2
	var respBundles []*dkg.ResponseBundle
	for _, n := range TestNodes {
		r, err := n.drand.ProcessDeals(deals)
		require.NoError(t, err)
		if r != nil {
			respBundles = append(respBundles, r)
		}
	}

	// Step 3
	var justifs []*dkg.JustificationBundle
	var results []*dkg.Result
	for _, n := range TestNodes {
		res, just, err := n.drand.ProcessResponses(respBundles)
		require.NoError(t, err)

		if res != nil {
			results = append(results, res)
		} else if just != nil {
			justifs = append(justifs, just)
		}
	}

	if len(justifs) > 0 {
		for _, n := range TestNodes {
			res, err := n.drand.ProcessJustifications(justifs)
			require.NoError(t, err)
			require.NotNil(t, res)
			results = append(results, res)
		}

		panic("implement")
	}

	require.NotEmpty(t, results)

	// print and collect individual share SKs
	var sks []bls3.SecretKey
	for i, res := range results {
		sk, err := resultToShareSecretKey(res)
		require.NoError(t, err)
		TestNodes[i].generatedShare = sk.Serialize()
		fmt.Printf("Index (%d): %x\n", i+1, sk.Serialize())
		sks = append(sks, *sk)
	}

	// reconstruct from shares and verify
	valSK := reconstructSK(t, sks)
	valPK, err := resultsToValidatorPK(results, TestSuite.G1().(dkg.Suite))
	require.NoError(t, err)
	require.EqualValues(t, valPK.Serialize(), valSK.GetPublicKey().Serialize())
	fmt.Printf("Validator SK: %x\nValidator PK: %x\n", valSK.Serialize(), valPK.Serialize())

	// encrypt shares
	var encryptedShares [][]byte
	for _, n := range TestNodes {
		cypher, err := n.EncryptShare()
		require.NoError(t, err)
		encryptedShares = append(encryptedShares, cypher)
	}

	// generate deposit data sigs
	root, _, err := GenerateETHDepositData(
		valPK.Serialize(),
		TestWithdrawalCredentials,
		GenesisForkVersion,
		DomainDeposit,
	)
	require.NoError(t, err)
	var depositPartialSigs []*bls3.Sign
	for _, n := range TestNodes {
		depositPartialSigs = append(depositPartialSigs, n.GetGenerateShare().SignByte(root))
	}

	// reconstruct deposit sig and verify
	sig, err := reconstructSignatures(map[int][]byte{
		1: depositPartialSigs[0].Serialize(),
		2: depositPartialSigs[1].Serialize(),
		3: depositPartialSigs[2].Serialize(),
		4: depositPartialSigs[3].Serialize(),
	})
	require.NoError(t, err)
	require.True(t, sig.VerifyByte(valPK, root))

	// construct output struct
	var output []*Output
	for i, _ := range TestNodes {
		output = append(output, &Output{
			Nonce:                       nonce,
			EncryptedShare:              encryptedShares[i],
			SharePK:                     TestNodes[i].GetGenerateShare().GetPublicKey().Serialize(),
			ValidatorPK:                 valPK.Serialize(),
			DepositDataPartialSignature: depositPartialSigs[i].Serialize(),
		})
	}
}
