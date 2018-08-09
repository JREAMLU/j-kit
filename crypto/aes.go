package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"math/rand"
	"time"

	"github.com/JREAMLU/j-kit/constant"
)

var (
	// ErrPaddingSize padding error
	ErrPaddingSize = errors.New("padding size error")
)

// AESEncrypter aes encrypt
func AESEncrypter(src string, encrypteKey string) (string, []byte, error) {
	key, err := hex.DecodeString(encrypteKey)
	if err != nil {
		return constant.EmptyStr, nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return constant.EmptyStr, nil, err
	}

	iv := randomIV()
	mode := cipher.NewCBCEncrypter(block, iv)

	content := []byte(src)
	content = PKCS7Padding(content, block.BlockSize())
	dst := make([]byte, len(content))
	mode.CryptBlocks(dst, content)

	ciphertext := base64.StdEncoding.EncodeToString(dst)

	return ciphertext, iv, nil
}

// AESDecrypter aes decrypt
func AESDecrypter(src string, encrypteKey string, iv []byte) ([]byte, error) {
	key, err := hex.DecodeString(encrypteKey)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	content, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	dst := make([]byte, len(content))
	mode.CryptBlocks(dst, content)

	return PKCS7UnPadding(dst, block.BlockSize())
}

// PKCS7Padding pkcs7 padding
func PKCS7Padding(src []byte, blockSize int) []byte {
	padding := blockSize - (len(src) % blockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(src, padText...)
}

// PKCS7UnPadding pkcs7 unpadding
func PKCS7UnPadding(src []byte, blockSize int) ([]byte, error) {
	length := len(src)
	unpadding := int(src[length-1])

	if unpadding >= length || unpadding > blockSize {
		return nil, ErrPaddingSize
	}

	return src[:(length - unpadding)], nil
}

func randomIV() (iv []byte) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 16; i++ {
		iv = append(iv, byte(rand.Intn(16)))
	}

	return iv
}