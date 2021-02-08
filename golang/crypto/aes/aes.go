package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"

	ccrypto "github.com/crossedbot/common/golang/crypto"
)

// NewKey returns a new GCM wrapped AES cipher block of a PBKDF2 extended key
func NewKey(key, salt []byte) (cipher.AEAD, []byte, error) {
	aead, err := AesGcmKey(ccrypto.ExtendKey(key, salt))
	if err != nil {
		return nil, nil, err
	}
	nonce, err := ccrypto.NewNonce(aead.NonceSize())
	if err != nil {
		return nil, nil, err
	}
	return aead, nonce, nil
}

// AesGcmKey returns a new GCM wrapped AES cipher block; keys should be 16, 24,
// or 32 bytes in length to select AES-128, AES-192, or AES-256 respectively
func AesGcmKey(key []byte) (cipher.AEAD, error) {
	switch len(key) {
	// check key length
	case 16, 24, 32:
	default:
		return nil, fmt.Errorf(
			"invalid key length (%d); %s",
			len(key),
			"accepted lengths are 16, 24, or 32 bytes",
		)
	}
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(c)
}
