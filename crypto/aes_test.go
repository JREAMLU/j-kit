package crypto

import (
	"encoding/base64"
	"encoding/hex"
	"math/rand"
	"testing"
	"time"

	"github.com/JREAMLU/j-kit/constant"
	"github.com/JREAMLU/j-kit/ext"
	. "github.com/smartystreets/goconvey/convey"
)

func TestEncrypt(t *testing.T) {
	key := "0F10F6CB2F5369C14D14FA07BAD302267901240CC8C845DD2C645FBD149A11C9"
	data := "123"

	Convey("encrypt test", t, func() {
		ciphertext, iv, err := AESEncrypter(data, key)
		t.Log(err)
		t.Log(iv)
		t.Log(ciphertext)
		So(err, ShouldBeNil)
		So(iv, ShouldNotBeNil)
		So(ciphertext, ShouldNotBeEmpty)

		raw, err := AESDecrypter(ciphertext, key, iv)
		t.Log(err)
		t.Log(string(raw))
		So(err, ShouldBeNil)
		So(string(raw), ShouldEqual, data)
	})
}

func TestCookie(t *testing.T) {
	// key, err := keyGen()
	// if err != nil {
	// 	panic(err)
	// }

	key := "0F10F6CB2F5369C14D14FA07BAD302267901240CC8C845DD2C645FBD149A11C9"
	data := `{"userID":10000}`
	// data := "10000000"

	Convey("cookie test", t, func() {
		cookie, err := cookieEncrpy(data, key)
		t.Log(err)
		t.Log(cookie)
		So(err, ShouldBeNil)
		So(cookie, ShouldNotBeEmpty)

		raw, err := cookieDecrypter(cookie, key)
		t.Log(err)
		t.Log(raw)
		So(err, ShouldBeNil)
		So(raw, ShouldEqual, data)
	})
}

func TestEDCookie(t *testing.T) {
	key, err := keyGen()
	if err != nil {
		panic(err)
	}
	vkey, err := keyGen()
	if err != nil {
		panic(err)
	}

	// key := "0F10F6CB2F5369C14D14FA07BAD302267901240CC8C845DD2C645FBD149A11C9"
	// vkey := "C985085862F161091EEEFE30F7DC9D62"
	data := `{"userID":10000}`

	Convey("encrypt decrypt cookie test", t, func() {
		cookie, err := EncryptCookie(data, key, vkey)
		t.Log(err)
		t.Log(cookie)
		So(err, ShouldBeNil)
		So(cookie, ShouldNotBeEmpty)

		raw, err := DecryptCookie(cookie, key, vkey)
		t.Log(err)
		t.Log(string(raw))
		So(err, ShouldBeNil)
		So(string(raw), ShouldEqual, data)
	})
}

func cookieEncrpy(data, key string) (string, error) {
	ciphertext, iv, err := AESEncrypter(data, key)
	if err != nil {
		return constant.EmptyStr, err
	}

	binCiphertext := hex.EncodeToString([]byte(ciphertext))

	base64IV := base64.StdEncoding.EncodeToString(iv)
	binIV := hex.EncodeToString([]byte(base64IV))

	return ext.StringSplice(binIV, binCiphertext), nil
}

func cookieDecrypter(cookie, key string) (string, error) {
	binIV := cookie[:48]
	base64IV, err := hex.DecodeString(binIV)
	if err != nil {
		return constant.EmptyStr, err
	}

	iv, err := base64.StdEncoding.DecodeString(string(base64IV))
	if err != nil {
		return constant.EmptyStr, err
	}

	binCiphertext := cookie[48:]
	ciphertext, err := hex.DecodeString(binCiphertext)
	if err != nil {
		return constant.EmptyStr, err
	}

	raw, err := AESDecrypter(string(ciphertext), key, iv)
	if err != nil {
		return constant.EmptyStr, err
	}

	return string(raw), nil
}

func keyGen() (string, error) {
	a, err := MD5(getRandomString(1024), true)
	if err != nil {
		return constant.EmptyStr, err
	}

	b, err := MD5(getRandomString(1024), true)
	if err != nil {
		return constant.EmptyStr, err
	}

	return ext.StringSplice(a, b), nil
}

func getRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}

	return string(result)
}
