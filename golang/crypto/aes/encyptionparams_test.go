package aes

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncrypt(t *testing.T) {
	plain := []byte("my secret message :3")
	key := []byte("averysecurekey")
	salt := []byte("somesalt")
	ae, nonce, err := NewKey(key, salt)
	require.Nil(t, err)
	expected := ae.Seal(nil, nonce, plain, nil)
	ep, err := NewEncryptionParamsWithValues(key, salt, nonce)
	require.Nil(t, err)
	actual := ep.Encrypt(plain)
	require.Equal(t, expected, actual)
}

func TestDecrypt(t *testing.T) {
	expected := []byte("my secret message :3")
	key := []byte("averysecurekey")
	salt := []byte("somesalt")
	ae, nonce, err := NewKey(key, salt)
	require.Nil(t, err)
	cipher := ae.Seal(nil, nonce, expected, nil)
	ep, err := NewEncryptionParamsWithValues(key, salt, nonce)
	require.Nil(t, err)
	actual, err := ep.Decrypt(cipher)
	require.Nil(t, err)
	require.Equal(t, expected, actual)
}

func TestNonce(t *testing.T) {
	key := []byte("averysecurekey")
	salt := []byte("somesalt")
	nonce := []byte("thisisanonce")
	ep, err := NewEncryptionParamsWithValues(key, salt, nonce)
	require.Nil(t, err)
	actual := ep.Nonce(nil)
	require.Equal(t, nonce, actual)
	nonce2 := []byte("thisisanothernonce")
	actual = ep.Nonce(nonce2)
	require.Equal(t, nonce2, actual)
}

func TestSalt(t *testing.T) {
	key := []byte("averysecurekey")
	salt := []byte("somesalt")
	nonce := []byte("thisisanonce")
	ep, err := NewEncryptionParamsWithValues(key, salt, nonce)
	require.Nil(t, err)
	actual := ep.Salt(nil)
	require.Equal(t, salt, actual)
	salt2 := []byte("someothersalt")
	actual = ep.Nonce(salt2)
	require.Equal(t, salt2, actual)

}
