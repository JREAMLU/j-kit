package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"math/rand"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/JREAMLU/j-kit/constant"
	"github.com/JREAMLU/j-kit/ext"
	"lukechampine.com/adiantum"
)

var (
	// ErrPaddingSize padding error
	ErrPaddingSize = errors.New("padding size error")
	// ErrSrcNotEmpty not empty
	ErrSrcNotEmpty = errors.New("src not empty")
	// ErrSrcMod err mod
	ErrSrcMod = errors.New("src mod no correct")
	// ErrHashCheck hash check
	ErrHashCheck = errors.New("hash check failed")
	// ErrPadding padding err
	ErrPadding = errors.New("padding err")
	// ErrPaddingEq padding diff
	ErrPaddingEq = errors.New("padding eq err")
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

// EncryptCookie aes encrypt cookie
func EncryptCookie(src string, encrypteKey string, validationKey string) (string, error) {
	key, err := hex.DecodeString(encrypteKey)
	if err != nil {
		return constant.EmptyStr, err
	}

	vkey, err := hex.DecodeString(validationKey)
	if err != nil {
		return constant.EmptyStr, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return constant.EmptyStr, err
	}

	iv := randomIV()
	ivLength := len(iv)

	padding := ivLength - (len(src) % ivLength)
	// repeat string() == chr
	src = ext.StringSplice(src, strings.Repeat(string(padding), padding))

	mode := cipher.NewCBCEncrypter(block, iv)

	content := []byte(src)
	content = PKCS7Padding(content, block.BlockSize())
	dst := make([]byte, len(content))
	mode.CryptBlocks(dst, content)

	// OPENSSL_RAW_DATA
	encrypted := string(dst)

	hashData := ext.StringSplice(string(iv), encrypted)
	hash, err := HMacSha256([]byte(hashData), string(vkey))
	if err != nil {
		return constant.EmptyStr, err
	}

	encryptedData := ext.StringSplice(hex.EncodeToString(iv), hex.EncodeToString([]byte(encrypted)), hash[:16])

	return strings.ToUpper(encryptedData), nil
}

// DecryptCookie aes decrypt cookie
func DecryptCookie(src string, encrypteKey string, validationKey string) ([]byte, error) {
	if len(src) <= 0 {
		return nil, ErrSrcNotEmpty
	}

	key, err := hex.DecodeString(encrypteKey)
	if err != nil {
		return nil, err
	}

	vkey, err := hex.DecodeString(validationKey)
	if err != nil {
		return nil, err
	}

	//check encryptData mod == 0
	if (len(src) % 2) != 0 {
		return nil, ErrSrcMod
	}

	binSrc, err := hex.DecodeString(src)
	if err != nil {
		return nil, err
	}
	src = string(binSrc)
	ivLength := len(randomIV())
	hashSize := 8
	hash := hex.EncodeToString([]byte(src[len(src)-hashSize:]))
	needHashData := src[:len(src)-hashSize]
	hashed, err := HMacSha256([]byte(needHashData), string(vkey))
	if err != nil {
		return nil, err
	}

	if hash != hashed[:16] {
		return nil, ErrHashCheck
	}

	iv := []byte(src[:ivLength])

	_src := src[ivLength : len(src)-hashSize]

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	dst := make([]byte, len(_src))
	mode.CryptBlocks(dst, []byte(_src))

	decryptedData, err := PKCS7UnPadding(dst, block.BlockSize())
	if err != nil {
		return nil, err
	}

	r, _ := utf8.DecodeRune(decryptedData[len(decryptedData)-1:])
	padding := int(r)

	if padding > len(decryptedData) {
		return nil, ErrPadding
	}

	if padding > strings.Count(string(decryptedData[len(decryptedData)-padding:]), string(padding)) {
		return nil, ErrHashCheck
	}

	return decryptedData[:len(decryptedData)-padding], nil
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

// EncryptAdiantumCookie aes encrypt adiantum cookie
func EncryptAdiantumCookie(src string, encrypteKey string, validationKey string) (string, error) {
	key, err := hex.DecodeString(encrypteKey)
	if err != nil {
		return constant.EmptyStr, err
	}

	vkey, err := hex.DecodeString(validationKey)
	if err != nil {
		return constant.EmptyStr, err
	}

	// tweak
	iv := randomIV()
	ivLength := len(iv)

	padding := ivLength - (len(src) % ivLength)
	// repeat string() == chr
	src = ext.StringSplice(src, strings.Repeat(string(padding), padding))
	content := []byte(src)

	// OPENSSL_RAW_DATA
	adCipher := adiantum.New20(key)
	encrypted := string(adCipher.Encrypt(content, iv))

	hashData := ext.StringSplice(string(iv), encrypted)
	hash, err := HMacSha256([]byte(hashData), string(vkey))
	if err != nil {
		return constant.EmptyStr, err
	}

	encryptedData := ext.StringSplice(hex.EncodeToString(iv), hex.EncodeToString([]byte(encrypted)), hash[:16])

	return strings.ToUpper(encryptedData), nil
}

// DecryptAdiantumCookie aes decrypt adiantum cookie
func DecryptAdiantumCookie(src string, encrypteKey string, validationKey string) ([]byte, error) {
	if len(src) <= 0 {
		return nil, ErrSrcNotEmpty
	}

	key, err := hex.DecodeString(encrypteKey)
	if err != nil {
		return nil, err
	}

	vkey, err := hex.DecodeString(validationKey)
	if err != nil {
		return nil, err
	}

	//check encryptData mod == 0
	if (len(src) % 2) != 0 {
		return nil, ErrSrcMod
	}

	binSrc, err := hex.DecodeString(src)
	if err != nil {
		return nil, err
	}
	src = string(binSrc)
	ivLength := len(randomIV())
	hashSize := 8
	hash := hex.EncodeToString([]byte(src[len(src)-hashSize:]))
	needHashData := src[:len(src)-hashSize]
	hashed, err := HMacSha256([]byte(needHashData), string(vkey))
	if err != nil {
		return nil, err
	}

	if hash != hashed[:16] {
		return nil, ErrHashCheck
	}

	// tweak
	iv := []byte(src[:ivLength])
	_src := src[ivLength : len(src)-hashSize]

	adCipher := adiantum.New20(key)
	decryptedData := adCipher.Decrypt([]byte(_src), iv)

	r, _ := utf8.DecodeRune(decryptedData[len(decryptedData)-1:])
	padding := int(r)

	if padding > len(decryptedData) {
		return nil, ErrPadding
	}

	if padding > strings.Count(string(decryptedData[len(decryptedData)-padding:]), string(padding)) {
		return nil, ErrHashCheck
	}

	return decryptedData[:len(decryptedData)-padding], nil
}

func randomIV() (iv []byte) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 16; i++ {
		iv = append(iv, byte(rand.Intn(16)))
	}

	return iv
}
