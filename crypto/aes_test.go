package crypto

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestEncrypt(t *testing.T) {
	key := "0F10F6CB2F5369C14D14FA07BAD302267901240CC8C845DD2C645FBD149A11C9"

	Convey("encrypt test", t, func() {
		ciphertext, iv, err := AESEncrypter("123", key)
		t.Log(err)
		t.Log(iv)
		t.Log(ciphertext)

		data, err := AESDecrypter(ciphertext, key, iv)
		t.Log(err)
		t.Log(string(data))
	})
}
