package dataconverter

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

var (
	ErrPadFailure = errors.New("unpad error. This could happen when incorrect encryption key is used")
	ErrIVFailure  = errors.New("failed random IV generation")
)

type AESEncryptionServiceV1 struct {
	CipherBlock cipher.Block
}

func newAESEncryptionServiceV1(opts Options) (*AESEncryptionServiceV1, error) {
	// must be 16, 24, 32 byte length
	// this is your encryption key
	// will fail to initialize if length requirements are not met
	cipherBlock, err := aes.NewCipher(opts.EncryptionKey)
	if err != nil {
		// likely invalid key length if errors here
		return nil, err
	}
	return &AESEncryptionServiceV1{
		CipherBlock: cipherBlock,
	}, nil
}

// Encrypt takes a byte array and returns an encrypted byte array
// as base64 encoded
func (a AESEncryptionServiceV1) Encrypt(unencryptedBytes []byte) ([]byte, error) {
	msg := pad(unencryptedBytes)
	cipherBytes := make([]byte, aes.BlockSize+len(msg))
	iv := cipherBytes[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		// this should never happen
		return nil, ErrIVFailure
	}
	cipher.
		NewCFBEncrypter(a.CipherBlock, iv).
		XORKeyStream(cipherBytes[aes.BlockSize:], msg)

	var encryptedBytes = make([]byte, base64.StdEncoding.EncodedLen(len(cipherBytes)))
	base64.StdEncoding.Encode(encryptedBytes ,cipherBytes)
	return encryptedBytes, nil
}

// Decrypt takes an encrypted base64 byte array then
// returns an unencrypted byte array if same key was used to encrypt it
func (a AESEncryptionServiceV1) Decrypt(encryptedBytes []byte) ([]byte, error) {
	if len(encryptedBytes) == 0 {
		return []byte(""), nil
	}
	decodeLen := (base64.StdEncoding.DecodedLen(len(encryptedBytes)) / aes.BlockSize) * aes.BlockSize
	decodedBytes := make([]byte, decodeLen)
	if _, err := base64.StdEncoding.Decode(decodedBytes, encryptedBytes); err != nil {
		return nil, err
	}
	iv := decodedBytes[:aes.BlockSize]
	msg := decodedBytes[aes.BlockSize:]
	cipher.
		NewCFBDecrypter(a.CipherBlock, iv).
		XORKeyStream(msg, msg)

	unpadMsg, err := unpad(msg)
	if err != nil {
		return nil, err
	}
	return unpadMsg, nil
}

// pad for AES encryption we need the message to be divisible by AES block size
func pad(src []byte) []byte {
	padding := aes.BlockSize - len(src)%aes.BlockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padText...)
}

// unpad remove padding we added before to allow AES encryption
func unpad(src []byte) ([]byte, error) {
	length := len(src)
	unpadding := int(src[length-1])
	if unpadding > length {
		return nil, ErrPadFailure
	}
	return src[:(length - unpadding)], nil
}


