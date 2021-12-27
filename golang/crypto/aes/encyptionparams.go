package aes

import (
	"crypto/cipher"

	ccrypto "github.com/crossedbot/common/golang/crypto"
)

// EncryptionParams represents an interface to perform symmetric encryption
type EncryptionParams interface {
	// Encrypt encrypts the given plain text and returns its cipher
	Encrypt(plain []byte) []byte

	// Decrypt decrypts the given cipher and returns its plain text
	Decrypt(cipher []byte) ([]byte, error)

	// Nonce (if given) sets the nonce and returns the current nonce
	Nonce(nonce []byte) []byte

	// Salt (if given) sets the salt and returns the current salt
	Salt(salt []byte) []byte
}

// encryptionParams represents implements the EncryptionParams interface using
// AES GCM encryption and tracks its state
type encryptionParams struct {
	cipher.AEAD
	nonce []byte
	salt  []byte
}

// NewEncryptionParams returns a new EncryptionParams, using the given key to
// generate the salt and nonce parameters
func NewEncryptionParams(key []byte) (EncryptionParams, error) {
	keyLen := len(key)
	b, err := ccrypto.GenerateRandomBytes(keyLen * 2)
	if err != nil {
		return nil, err
	}
	salt := b[:keyLen]
	ae, nonce, err := NewKey(key, salt)
	if err != nil {
		return nil, err
	}
	return &encryptionParams{
		AEAD:  ae,
		nonce: nonce,
		salt:  salt,
	}, nil
}

// NewEncryptionParamsWithValues returns a new EncryptionParams using the given
// parameter values
func NewEncryptionParamsWithValues(key, salt, nonce []byte) (EncryptionParams, error) {
	var err error
	encParams := encryptionParams{nonce: nonce, salt: salt}
	encParams.AEAD, err = AesGcmKey(ccrypto.ExtendKey(key, salt))
	if err != nil {
		return nil, err
	}
	return &encParams, nil
}

// Encrypts encrypts the given plain text
func (ep *encryptionParams) Encrypt(plain []byte) []byte {
	return ep.AEAD.Seal(nil, ep.nonce, plain, nil)
}

// Decrypt decrypts the given cipher text
func (ep *encryptionParams) Decrypt(cipher []byte) ([]byte, error) {
	return ep.AEAD.Open(nil, ep.nonce, cipher, nil)
}

// Nonce sets and returns the current nonce
func (ep *encryptionParams) Nonce(nonce []byte) []byte {
	if nonce != nil && len(nonce) > 0 {
		ep.nonce = nonce
	}
	return ep.nonce
}

// Salt sets and returns the current salt
func (ep *encryptionParams) Salt(salt []byte) []byte {
	if salt != nil && len(salt) > 0 {
		ep.salt = salt
	}
	return ep.salt
}
