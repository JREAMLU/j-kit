package sign

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/JREAMLU/j-core/constant"
	"github.com/JREAMLU/j-core/crypto"
	"github.com/JREAMLU/j-core/ext"
)

//GenerateSign 生成签名 参数key全部按键值排序     ToUpper(md5(sha1(base64(urlencode(SecretKey1Value1Key2Value2SecretTime)))))
//strtoupper( md5 ( sha1( base64_encode( urlencode( secret_key . static::serialize( request_data ) . secret_key . request_time ) ) ) ) )
func GenerateSign(requestData []byte, requestTime int64, secret string) (string, error) {
	var rdata map[string]interface{}
	err := json.Unmarshal(requestData, &rdata)
	if err != nil {
		return constant.EmptyStr, err
	}

	str := serialize(rdata)
	serial := ext.StringSplice(secret, str.(string), secret, strconv.FormatInt(int64(requestTime), 10))
	urlencodeSerial := url.QueryEscape(serial)
	urlencodeBase64Serial := base64.StdEncoding.EncodeToString([]byte(urlencodeSerial))
	sign, err := crypto.Sha1(urlencodeBase64Serial)
	if err != nil {
		return constant.EmptyStr, err
	}
	sign, err = crypto.MD5(sign)
	if err != nil {
		return constant.EmptyStr, err
	}

	return strings.ToUpper(sign), nil
}

// Serialize 序列化 && 递归ksort
func serialize(data interface{}) interface{} {
	var buffer bytes.Buffer
	switch reflect.TypeOf(data).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(data)
		for i := 0; i < s.Len(); i++ {
			serial := serialize(s.Index(i).Interface())
			if reflect.TypeOf(serial).Kind() == reflect.Float64 {
				serial = strconv.Itoa(int(serial.(float64)))
			}
			buffer.WriteString(strconv.Itoa(i))
			buffer.WriteString(serial.(string))
		}
		return buffer.String()
	case reflect.Map:
		s := reflect.ValueOf(data)
		keys := s.MapKeys()
		//ksort
		var sortedKeys []string
		for _, key := range keys {
			sortedKeys = append(sortedKeys, key.Interface().(string))
		}
		sort.Strings(sortedKeys)
		for _, key := range sortedKeys {
			serial := serialize(s.MapIndex(reflect.ValueOf(key)).Interface())
			if reflect.TypeOf(serial).Kind() == reflect.Float64 {
				serial = strconv.Itoa(int(serial.(float64)))
			}
			buffer.WriteString(key)
			buffer.WriteString(serial.(string))
		}
		return buffer.String()
	}

	return data
}

// ValidSign 签名验证
func ValidSign(requestData []byte, sign string, timestamp int64, secret string, expire int64) error {
	var rdata map[string]interface{}
	json.Unmarshal(requestData, &rdata)

	jsonData, err := json.Marshal(rdata)
	if err != nil {
		return err
	}

	signed, err := GenerateSign(jsonData, timestamp, secret)
	if err != nil {
		return err
	}

	if sign != signed {
		return errors.New("INVALID SIGNATURE")
	}

	if diff := time.Now().Unix() - timestamp; diff > expire {
		return errors.New("SIGNATURE TIMEEXPIRED")
	}

	return nil
}

// PHP
/*
class client {

	public function generateMKII($request_data = [], $request_time = "", $secret_key = "") {
		return strtoupper( md5 ( sha1( base64_encode( urlencode( $secret_key . static::serialize( $request_data ) . $secret_key . $request_time ) ) ) ) );
	}

	public static function serialize($data) {
		if (is_array($data)) {
			ksort($data);
			$str = "";
			foreach ($data as $key => $value) {
				$str = sprintf('%s%s%s', $str, $key, static::serialize($value));
			}
			return $str;
		} else {
			return $data;
		}
	}

}
*/
