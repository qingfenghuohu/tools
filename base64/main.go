package base64

import (
	"encoding/base64"
	"strings"
)

func SafeEncode(source []byte) string {
	bytearr := base64.StdEncoding.EncodeToString(source)
	safeurl := strings.Replace(string(bytearr), "/", "_", -1)
	safeurl = strings.Replace(safeurl, "+", "-", -1)
	safeurl = strings.Replace(safeurl, "=", "", -1)
	return safeurl
}

func Encode(source string) string {
	bytearr := base64.StdEncoding.EncodeToString([]byte(source))
	safeurl := strings.Replace(string(bytearr), "/", "_", -1)
	safeurl = strings.Replace(safeurl, "+", "-", -1)
	safeurl = strings.Replace(safeurl, "=", "", -1)
	return safeurl
}

func Decode(source string) string {
	result, _ := base64.StdEncoding.DecodeString(source)
	return string(result)
}
