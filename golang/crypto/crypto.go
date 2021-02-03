package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"strconv"
	"time"

	"golang.org/x/crypto/pbkdf2"
)

const (
	// Encryption Constants
	AuthKeyIdSize   = 8
	KdfIterations   = 4096
	ExtendedKeySize = 32
)

// generateRandomBytes returns n number of random bytes
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	return b, nil
}

// generateRandomString returns a random string of n length
func GenerateRandomString(n int) (string, error) {
	b, err := GenerateRandomBytes(n * 2)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b)[:n], nil
}

// NewNonce returns a new nonce for the given size
func NewNonce(sz int) ([]byte, error) {
	now := time.Now().Unix()
	h := sha256.New()
	str, err := GenerateRandomString(h.Size())
	if err != nil {
		return nil, err
	}
	io.WriteString(h, strconv.FormatInt(now, 10))
	io.WriteString(h, str)
	return h.Sum(nil)[:sz], nil
}

// ExtendKey returns an extended key using the PBKDF2 function
func ExtendKey(key, salt []byte) []byte {
	return pbkdf2.Key(
		key,             // password
		salt,            // salt
		KdfIterations,   // iterations
		ExtendedKeySize, // key size
		sha256.New,      // hash function
	)
}

// KeyId returns a AuthKeyIdSize long ID for the given key
func KeyId(key []byte) []byte {
	sum := sha256.Sum256(key)
	return sum[:AuthKeyIdSize]
}
