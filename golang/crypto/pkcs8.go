package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

var (
	// Errors
	ErrKeyNotFound          = errors.New("Key not found")
	ErrNotEncryptedPEMBlock = errors.New("PEM block is not encrypted")
	ErrUnsupportedAlgorithm = errors.New("Unsupported encryption algorithm")
	ErrUnsupportedFormat    = errors.New("Unsupported encryption format, expecting 'PBES2'")
	ErrUnsupportedKDF       = errors.New("Unsupported KDF, expecting PBKD2")
)

var (
	// ASN.1 Formats
	// RFC8018 Appendix A, RFC8018 Appendix C
	OidRSADI                 = asn1.ObjectIdentifier{1, 2, 840, 113549}
	OidPKCS                  = asn1.ObjectIdentifier{1, 2, 840, 113549, 1}
	OidPKCS5                 = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 5}
	OidPBEWithMD2AndDES_CBC  = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 5, 1}
	OidPBEWithMD5AndDES_CBC  = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 5, 3}
	OidPBEWithMD2AndRC2_CBC  = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 5, 4}
	OidPBEWithMD5AndRC2_CBC  = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 5, 6}
	OidPBEWithSHA1AndDES_CBC = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 5, 10}
	OidPBEWithSHA1AndRC2_CBC = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 5, 11}
	OidPBKDF2                = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 5, 12}
	OidPBES2                 = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 5, 13}
	OidPBMAC1                = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 5, 14}

	// Supported KFDs and Encryption Schemes
	// RFC8018 Appendix A, RFC8018 Appendix C
	OidDigestAlgorithm    = asn1.ObjectIdentifier{1, 2, 840, 113549, 2}
	OidHMACWithSHA1       = asn1.ObjectIdentifier{1, 2, 840, 113549, 2, 7}
	OidHMACWithSHA224     = asn1.ObjectIdentifier{1, 2, 840, 113549, 2, 8}
	OidHMACWithSHA256     = asn1.ObjectIdentifier{1, 2, 840, 113549, 2, 9}
	OidHMACWithSHA384     = asn1.ObjectIdentifier{1, 2, 840, 113549, 2, 10}
	OidHMACWithSHA512     = asn1.ObjectIdentifier{1, 2, 840, 113549, 2, 11}
	OidHMACWithSHA512_224 = asn1.ObjectIdentifier{1, 2, 840, 113549, 2, 12}
	OidHMACWithSHA512_256 = asn1.ObjectIdentifier{1, 2, 840, 113549, 2, 13}

	OidEncryptionAlgorithm = asn1.ObjectIdentifier{1, 2, 840, 113549, 3}
	OidRC2CBC              = asn1.ObjectIdentifier{1, 2, 840, 113549, 3, 2}
	OidDES_EDE3_CBC        = asn1.ObjectIdentifier{1, 2, 840, 113549, 3, 7}
	OidRC2_CBC_PAD         = asn1.ObjectIdentifier{1, 2, 840, 113549, 3, 9}

	OidOIW    = asn1.ObjectIdentifier{1, 3, 14}
	OidDESCBC = asn1.ObjectIdentifier{1, 3, 14, 3, 2, 7}

	OidNistAlgorithms = asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4}
	OidAES            = asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 1}
	OidAES128_CBC_PAD = asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 1, 2}
	OidAES192_CBC_PAD = asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 1, 22}
	OidAES256_CBC_PAD = asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 1, 42}
)

// RFC8018 Appendix A.2 PBKDF2-PRFs
type PBKDF2PRFs struct {
	Algorithm asn1.ObjectIdentifier
	NullID    asn1.RawValue
}

// RFC8018 Appendix A.2 PKBKDF2-params
type PBKDF2Params struct {
	Salt           []byte
	IterationCount int
	PRF            PBKDF2PRFs `asn1:"optional"`
}

// RFC8018 Appendix A.4 PBES2-KDFs
type PBES2KDFs struct {
	Algorithm  asn1.ObjectIdentifier
	Parameters PBKDF2Params
}

// RFC8018 Appendix B.2 - B.2.2 DES-EDE3-CBC-Pad (Probably)
type PBES2Encs struct {
	Algorithm asn1.ObjectIdentifier
	IV        []byte
}

// RFC8018 Appendix A.4 PBES2-params
type PBES2Params struct {
	KeyDerivationFunc PBES2KDFs
	EncryptionScheme  PBES2Encs
}

// RFC5280 Section 4.1.1.2 AlgorithmIdentifer
type AlgorithmIdentifier struct {
	Algorithm  asn1.ObjectIdentifier
	Parameters PBES2Params
}

// RFC5208 Section 6 EncryptedPrivateKeyInfo
type EncryptedPrivateKeyInfo struct {
	EncryptionAlgorithm AlgorithmIdentifier
	EncryptedData       []byte
}

// DecryptPrivateKey returns the decrypted PEM block for the given PEM encoded
// private key and passphrase.
func DecryptPrivateKey(key, password []byte) ([]byte, error) {
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, ErrKeyNotFound
	}
	der, err := DecryptPEMBlock(block, password)
	if err != nil {
		return nil, err
	}
	blockType := strings.ReplaceAll(block.Type, "ENCRYPTED ", "")
	return pem.EncodeToMemory(&pem.Block{Type: blockType, Bytes: der}), nil
}

// DecryptPEMBlock returns the decrypted PEM block using the given passphrase.
func DecryptPEMBlock(b *pem.Block, password []byte) ([]byte, error) {
	if b.Headers["Proc-Type"] == "4,ENCRYPTED" {
		return x509.DecryptPEMBlock(b, password)
	}
	if b.Type == "ENCRYPTED PRIVATE KEY" {
		return DecryptPKCS8Key(b.Bytes, password)
	}
	return nil, ErrNotEncryptedPEMBlock
}

// DecryptPKCS8Key decrypts the given PKCS#8 formatted DER encoded ASN.1
// structure, and returns it decrypted using the given passphrase.
func DecryptPKCS8Key(data, password []byte) ([]byte, error) {
	var pki EncryptedPrivateKeyInfo
	_, err := asn1.Unmarshal(data, &pki)
	if err != nil {
		return nil, err
	}
	if !pki.EncryptionAlgorithm.Algorithm.Equal(OidPBES2) {
		return nil, ErrUnsupportedFormat
	}
	if !pki.EncryptionAlgorithm.Parameters.KeyDerivationFunc.Algorithm.Equal(OidPBKDF2) {
		return nil, ErrUnsupportedKDF
	}
	kdf := pki.EncryptionAlgorithm.Parameters.KeyDerivationFunc
	scheme := pki.EncryptionAlgorithm.Parameters.EncryptionScheme
	salt := kdf.Parameters.Salt
	iter := kdf.Parameters.IterationCount
	iv := scheme.IV
	prf := kdf.Parameters.PRF
	hashFn := sha1.New
	switch {
	case prf.Algorithm.Equal(OidHMACWithSHA224):
		hashFn = sha256.New224
	case prf.Algorithm.Equal(OidHMACWithSHA256):
		hashFn = sha256.New
	case prf.Algorithm.Equal(OidHMACWithSHA384):
		hashFn = sha512.New384
	case prf.Algorithm.Equal(OidHMACWithSHA512):
		hashFn = sha512.New
	case prf.Algorithm.Equal(OidHMACWithSHA512_224):
		hashFn = sha512.New512_224
	case prf.Algorithm.Equal(OidHMACWithSHA512_256):
		hashFn = sha512.New512_256
	}
	var keyLen int
	var newCipher func([]byte) (cipher.Block, error)
	switch {
	case scheme.Algorithm.Equal(OidDES_EDE3_CBC):
		keyLen = 8
		newCipher = des.NewCipher
	case scheme.Algorithm.Equal(OidDESCBC):
		keyLen = 24
		newCipher = des.NewTripleDESCipher
	case scheme.Algorithm.Equal(OidAES128_CBC_PAD):
		keyLen = 16
		newCipher = aes.NewCipher
	case scheme.Algorithm.Equal(OidAES192_CBC_PAD):
		keyLen = 16
		newCipher = aes.NewCipher
	case scheme.Algorithm.Equal(OidAES256_CBC_PAD):
		keyLen = 32
		newCipher = aes.NewCipher
	default:
		return nil, ErrUnsupportedAlgorithm
	}
	key := pbkdf2.Key(password, salt, iter, keyLen, hashFn)
	block, err := newCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	blockMode.CryptBlocks(pki.EncryptedData, pki.EncryptedData)
	return pki.EncryptedData, nil
}
