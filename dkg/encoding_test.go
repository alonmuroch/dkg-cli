package dkg

import (
	"encoding/hex"
	"fmt"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/stretchr/testify/require"
	"testing"
)

func skFromString(t *testing.T, s string) bls.SecretKey {
	ret := bls.SecretKey{}
	byts, _ := hex.DecodeString(s)
	require.NoError(t, ret.Deserialize(byts))
	return ret
}

func blsID(t *testing.T, id int) bls.ID {
	blsID := bls.ID{}
	require.NoError(t, blsID.SetDecString(fmt.Sprintf("%d", id)))
	return blsID
}

// tested against the below encoding from drand
/*
// taken from https://github.com/alonmuroch/kyber/blob/bls-dkg/share/dkg/dkg_test.go#L126-L173
	bytssk, _ := share.V.MarshalBinary()
	sk := &bls2.SecretKey{}
	require.NoError(t, sk.Deserialize(bytssk))
	fmt.Printf("sk: %s\n", hex.EncodeToString(sk.Serialize()))

	byts, _ := pubShare.V.MarshalBinary()
	fmt.Printf("pubshare: %x\n", byts)
	pk := &bls2.PublicKey{}
	require.NoError(t, pk.Deserialize(byts))
	fmt.Printf("pk hex: %s\n", pk.GetHexString())
*/
func TestEncodeFromDrand(t *testing.T) {
	_ = bls.Init(bls.BLS12_381)
	_ = bls.SetETHmode(bls.EthModeDraft07)

	sk1 := skFromString(t, "3139da31de82f972f5422c1aa6b6596711084edc13ed8efb4afbac3d4f050beb")
	sk2 := skFromString(t, "5b7bb518af1b91923ada8d2347a4a6b8347a272d50c0d1f076458a7720c6caa4")
	sk3 := skFromString(t, "1c4cc57a448c8514f4a922724161a77db8164c6caaf36be290814eabcbadfc03")
	sk4 := skFromString(t, "5b8859fcf210ce8b89219c17a7310bc2435806a0228214cf99aef8d94fbaa00a")
	validatorSK := skFromString(t, "1174dc18fc6039ff5719d7606838978fa17e677bf477ff020ea3b3fd5668bfd9")
	//validatorPK := "8eac72c83c7e416fab9cc1933c5b73702d4fbf83819738d12af33331ed85afb3df99f43502697861919e7c40daa2e93d"

	s := bls.SecretKey{}
	require.NoError(t, s.Recover(
		[]bls.SecretKey{sk1, sk2, sk3, sk4},
		[]bls.ID{blsID(t, 1), blsID(t, 2), blsID(t, 3), blsID(t, 4)},
	))

	require.EqualValues(t, validatorSK.GetHexString(), s.GetHexString())
}
