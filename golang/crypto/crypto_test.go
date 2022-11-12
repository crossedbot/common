package crypto

import (
	"crypto/x509"
	"encoding/pem"
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

func TestFingerprint(t *testing.T) {
	expected := "80E904715F82F6178C380D7F30D9388072B48E2E5C2BD769926726D280F1B575"
	block, _ := pem.Decode([]byte(`-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAtyYR8CCFssjvobNaVaQO
+qlte/usoVIrCyvWBX8Pod0Iq4HOma/nwBxP3k31OqF6pQ92BZq2fL8SBvPAceJR
rzL9TnBULkUS7wTbWIYuGJAnBXQ4+04NPMXFNEq0X6R2lyW58OAgduoWpcloYUGS
BVQuw9wnxSl8M3SQ76WISfTr/urNOtDz8E/zE3JFskkUzDzDnvXPqjResChE1flr
e5NLN3LDiSfB2OyjbYvcEFC+kwJcsseJhhUW+8EeoxFn8EwIUbWqIXFMGln/Rq7S
wwYb/C9yKKiieTa5+s8RW50Hy10gQXvI6Rl04hvqeVFRmRrb4Ga/aewFvksoL/vT
wNxVVGcw+wl0G0f4kqPxeftjy4pItKNBpeLKJebP5IcBZsx1wQTSgVoLxebIEw84
nog4YZoOqoB77vL24N8H6t6Xu8MYsa7spDkP9NIQ3GO7P/nupSMfIzZun5Nxaz6g
K38f5KIO+h3QnJejhj6i59htkeNZcGtLoBAb9OpTZj5X+1DEB1t/FTsfUKrE5Hpr
FYSXygYT3wxekedkalOWam15DHrwCKs6UU2DziUQR9VghjO7ILBNQ4cXZD6qEc/C
Q4Fno4es4QR4hEQwZSjw6DyiJRvOyDC+Ree9K+tJeV1wZawcfYSbpHdOuyu+wwqR
1ZwRFFjl+qJTV1y5MMx2DFkCAwEAAQ==
-----END PUBLIC KEY-----
`))
	require.NotNil(t, block)
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	require.Nil(t, err)
	actual := Fingerprint(pubKey)
	require.Equal(t, expected, actual)
}

func TestGetKeyId(t *testing.T) {
	key := []byte("helloworld")
	expected := []byte{
		0x93, 0x6a, 0x18, 0x5c, 0xaa, 0xa2, 0x66, 0xbb,
		0x44, 0x12, 0xbb, 0x6f, 0x8f, 0x8f, 0x07, 0xaf,
	}
	require.Equal(t, expected, KeyId(key))
}
