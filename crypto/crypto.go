package crypto

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"hash"
	"io"
	"strings"
)

//MD5 md5
func MD5(src string, isCaps ...bool) (string, error) {
	if len(isCaps) > 0 && isCaps[0] == true {
		s, err := _hash(src, md5.New())
		return strings.ToUpper(s), err
	}

	return _hash(src, md5.New())
}

//Sha1 sha1
func Sha1(src string, isCaps ...bool) (string, error) {
	if len(isCaps) > 0 && isCaps[0] == true {
		s, err := _hash(src, sha1.New())
		return strings.ToUpper(s), err
	}

	return _hash(src, sha1.New())
}

//HMacMD5 hmac-md5
func HMacMD5(src []byte, key string, isCaps ...bool) (string, error) {
	if len(isCaps) > 0 && isCaps[0] == true {
		s, err := _hmac(src, key, md5.New)
		return strings.ToUpper(s), err
	}

	return _hmac(src, key, md5.New)
}

//HMacSha1 hmac-sha1
func HMacSha1(src []byte, key string, isCaps ...bool) (string, error) {
	if len(isCaps) > 0 && isCaps[0] == true {
		s, err := _hmac(src, key, sha1.New)
		return strings.ToUpper(s), err
	}

	//hmac ,use sha1
	return _hmac(src, key, sha1.New)
}

func _hmac(src []byte, key string, h func() hash.Hash) (string, error) {
	mac := hmac.New(h, []byte(key))
	if _, err := mac.Write(src); err != nil {
		return "", err
	}

	return hex.EncodeToString(mac.Sum(nil)), nil
}

func _hash(src string, h hash.Hash) (string, error) {
	if _, err := io.WriteString(h, src); err != nil {
		return "", err
	}

	sig := hex.EncodeToString(h.Sum(nil))

	return sig, nil
}
