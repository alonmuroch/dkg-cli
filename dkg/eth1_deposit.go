package dkg

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"
	"github.com/pkg/errors"
)

// MaxEffectiveBalanceInGwei is the max effective balance
const MaxEffectiveBalanceInGwei uint64 = 32000000000

// BLSWithdrawalPrefixByte is the BLS withdrawal prefix
const BLSWithdrawalPrefixByte = byte(0)

var GenesisValidatorsRoot = phase0.Root{}
var GenesisForkVersion = phase0.Version{0, 0, 0, 0}
var DomainDeposit = [4]byte{0x03, 0x00, 0x00, 0x00}

// GenerateETHDepositData returns un-signed deposit data and deposit data root for signature
func GenerateETHDepositData(
	validatorPK, withdrawalCredentials []byte,
	fork phase0.Version,
	domain phase0.DomainType) ([]byte, *phase0.DepositData, error) {
	pk := phase0.BLSPubKey{}
	copy(pk[:], validatorPK)

	ret := &phase0.DepositMessage{
		PublicKey:             pk,
		WithdrawalCredentials: withdrawalCredentials,
		Amount:                phase0.Gwei(MaxEffectiveBalanceInGwei),
	}

	domainR, err := ComputeETHDomain(domain, fork, GenesisValidatorsRoot)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not compute deposit domain")
	}
	signingRoot, err := ComputeETHSigningRoot(ret, domainR)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not compute deposit signing root")
	}
	return signingRoot[:], &phase0.DepositData{
		PublicKey:             pk,
		WithdrawalCredentials: withdrawalCredentials,
		Amount:                phase0.Gwei(MaxEffectiveBalanceInGwei),
	}, nil
}

// ComputeETHDomain returns computed domain
func ComputeETHDomain(domain phase0.DomainType, fork phase0.Version, genesisValidatorRoot phase0.Root) (phase0.Domain, error) {
	ret := phase0.Domain{}
	copy(ret[0:4], domain[:])

	forkData := phase0.ForkData{
		CurrentVersion:        fork,
		GenesisValidatorsRoot: genesisValidatorRoot,
	}
	forkDataRoot, err := forkData.HashTreeRoot()
	if err != nil {
		return ret, err
	}
	copy(ret[4:32], forkDataRoot[0:28])
	return ret, nil
}

func ComputeETHSigningRoot(obj ssz.HashRoot, domain phase0.Domain) (phase0.Root, error) {
	root, err := obj.HashTreeRoot()
	if err != nil {
		return phase0.Root{}, err
	}
	signingContainer := phase0.SigningData{
		ObjectRoot: root,
		Domain:     domain,
	}
	return signingContainer.HashTreeRoot()
}
