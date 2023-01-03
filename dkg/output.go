package dkg

type Output struct {
	Nonce                       []byte
	EncryptedShare              []byte
	SharePK                     []byte
	ValidatorPK                 []byte
	DepositDataPartialSignature []byte
}
