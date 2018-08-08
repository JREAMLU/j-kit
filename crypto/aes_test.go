package crypto

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestEncrypt(t *testing.T) {
	Convey("encrypt test", t, func() {
		ciphertext, token, err := AESEncrypter("123")
		h, _ := HMacSha1([]byte(ciphertext), "abc", true)
		fmt.Println("++++++++++++: ", h)
		t.Log(err)
		t.Log(ciphertext)
		t.Log(token)

		data, err := AESDecrypter(ciphertext)
		t.Log(err)
		t.Log(string(data))
	})
}
