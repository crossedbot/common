package aes

import (
	"testing"

	"github.com/stretchr/testify/require"

	ccrypto "github.com/crossedbot/common/golang/crypto"
)

func TestAesGcmKey(t *testing.T) {
	key, salt := []byte("helloworld"), []byte("saltsalt")
	extKey := ccrypto.ExtendKey(key, salt)
	aead, err := AesGcmKey(extKey)
	require.Nil(t, err)
	require.NotNil(t, aead)
	require.NotZero(t, aead.NonceSize())
}

func TestNewKey(t *testing.T) {
	key, salt := []byte("helloworld"), []byte("saltsalt")
	aead, nonce, err := NewKey(key, salt)
	require.Nil(t, err)
	require.NotNil(t, aead)
	require.NotZero(t, aead.NonceSize())
	require.Equal(t, aead.NonceSize(), len(nonce))
}
