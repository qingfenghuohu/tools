package json

import (
	jsoniter "github.com/json-iterator/go"
	"strings"
)

func Encode(data interface{}) string {
	result, _ := jsoniter.Marshal(data)
	return string(result)
}

func Decode(data string) map[string]interface{} {
	var result map[string]interface{}
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	reader := strings.NewReader(string(data))
	decoder := json.NewDecoder(reader)
	decoder.Decode(&result)
	return result
}
