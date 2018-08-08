package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"math/rand"
	"strings"
)

const encrypteKey = "0F10F6CB2F5369C14D14FA07BAD302267901240CC8C845DD2C645FBD149A11C9"
const validationKey = "C985085862F161091EEEFE30F7DC9D62"

var (
	// DefaultAESIV default aes iv
	DefaultAESIV = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
)

// AESEncrypter aes encrypt
func AESEncrypter(src string) (string, string, error) {
	var keytoken []byte

	for i := 0; i < 16; i++ {
		keytoken = append(keytoken, byte(rand.Intn(255)))
	}

	key, err := hex.DecodeString(encrypteKey)
	if err != nil {
		return "", "", err
	}

	iv := DefaultAESIV

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)

	content := []byte(src)
	content = PKCS7Padding(content, block.BlockSize())
	dst := make([]byte, len(content))
	mode.CryptBlocks(dst, content)

	ciphertext := base64.StdEncoding.EncodeToString(dst)
	token := strings.ToUpper(hex.EncodeToString(keytoken))

	return ciphertext, token, nil
}

// AESDecrypter aes decrypt
func AESDecrypter(src string) ([]byte, error) {
	key, err := hex.DecodeString(encrypteKey)
	if err != nil {
		return nil, err
	}

	iv := DefaultAESIV

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

	return PKCS7UnPadding(dst), nil
}

// PKCS7Padding pkcs7 padding
func PKCS7Padding(src []byte, blockSize int) []byte {
	padding := blockSize - (len(src) % blockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(src, padText...)
}

// PKCS7UnPadding pkcs7 unpadding
func PKCS7UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])

	return src[:(length - unpadding)]
}

/*
func ExampleNewCBCEncrypter() {
	// Load your secret key from a safe place and reuse it across multiple
	// NewCipher calls. (Obviously don't use this example key for anything
	// real.) If you want to convert a passphrase to a key, use a suitable
	// package like bcrypt or scrypt.
	key, _ := hex.DecodeString("6368616e676520746869732070617373")
	plaintext := []byte("aaaaaaaaaaaaaaaa")

	// CBC mode works on blocks so plaintexts may need to be padded to the
	// next whole block. For an example of such padding, see
	// https://tools.ietf.org/html/rfc5246#section-6.2.3.2. Here we'll
	// assume that the plaintext is already of the correct length.
	if len(plaintext)%aes.BlockSize != 0 {
		panic("plaintext is not a multiple of the block size")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	// It's important to remember that ciphertexts must be authenticated
	// (i.e. by using crypto/hmac) as well as being encrypted in order to
	// be secure.

	fmt.Printf("%x\n", ciphertext)
}

func ExampleNewCBCDecrypter() {
	// Load your secret key from a safe place and reuse it across multiple
	// NewCipher calls. (Obviously don't use this example key for anything
	// real.) If you want to convert a passphrase to a key, use a suitable
	// package like bcrypt or scrypt.
	key, _ := hex.DecodeString("6368616e676520746869732070617373")
	ciphertext, _ := hex.DecodeString("42867bd1fad7d12661e62ecc1c3376cf0e0e9355b362ea190167a10e5c934541")

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// CBC mode always works in whole blocks.
	if len(ciphertext)%aes.BlockSize != 0 {
		panic("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(ciphertext, ciphertext)

	// If the original plaintext lengths are not a multiple of the block
	// size, padding would have to be added when encrypting, which would be
	// removed at this point. For an example, see
	// https://tools.ietf.org/html/rfc5246#section-6.2.3.2. However, it's
	// critical to note that ciphertexts must be authenticated (i.e. by
	// using crypto/hmac) before being decrypted in order to avoid creating
	// a padding oracle.

	fmt.Printf("%s\n", ciphertext)
	// Output: exampleplaintext
}
*/
