package dkg

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGenerateKey(t *testing.T) {
	sk, pk, err := GenerateKey()
	require.NoError(t, err)
	fmt.Printf("sk: %s\npk: %s\n", string(sk), string(pk))
}
