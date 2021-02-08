package crypto

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateRandomBytes(t *testing.T) {
	length := 12
	b1, err := GenerateRandomBytes(length)
	require.Nil(t, err)
	require.Equal(t, length, len(b1))
	b2, err := GenerateRandomBytes(length)
	require.Nil(t, err)
	require.NotEqual(t, b1, b2)
}

func TestGenerateRandomString(t *testing.T) {
	length := 12
	b1, err := GenerateRandomString(length)
	require.Nil(t, err)
	require.Equal(t, length, len(b1))
	b2, err := GenerateRandomString(length)
	require.Nil(t, err)
	require.NotEqual(t, b1, b2)
}

func TestNewNonce(t *testing.T) {
	length := 12
	n1, err := NewNonce(length)
	require.Nil(t, err)
	require.Equal(t, length, len(n1))
	n2, err := NewNonce(length)
	require.Nil(t, err)
	require.NotEqual(t, n1, n2)
}

func TestExtendKey(t *testing.T) {
	key, salt := []byte("helloworld"), []byte("saltsalt")
	b := ExtendKey(key, salt)
	require.Equal(t, ExtendedKeySize, len(b))
	require.Equal(t, b, ExtendKey(key, salt))
}

func TestGetKeyId(t *testing.T) {
	key := []byte("helloworld")
	expected := []byte{0x93, 0x6a, 0x18, 0x5c, 0xaa, 0xa2, 0x66, 0xbb}
	require.Equal(t, expected, KeyId(key))
}
