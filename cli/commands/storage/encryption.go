package storage

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha512"
	"github.com/pkg/errors"
	"golang.org/x/crypto/pbkdf2"
	"io"
)

func getAESKeyFromPassword(keyStr string) ([]byte, error) {
	return pbkdf2.Key([]byte(keyStr), []byte{}, 100000, 32, sha512.New), nil
}

// encrypt string to base64 crypto using AES
func encrypt(keyStr string, plaintext []byte) ([]byte, error) {
	key, err := getAESKeyFromPassword(keyStr)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// convert to base64
	return ciphertext, nil
}

// decrypt from base64 to decrypted string
func decrypt(keyStr string, ciphertext []byte) ([]byte, error) {
	key, err := getAESKeyFromPassword(keyStr)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, nil
}
