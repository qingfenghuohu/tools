package str

import (
	"encoding/json"
	"strconv"
	"strings"
	"unicode/utf8"
)

func Obj2Str(value interface{}) string {
	var key string
	if value == nil {
		return key
	}
	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	case bool:
		key = strconv.FormatBool(value.(bool))
	default:
		newValue, _ := json.Marshal(value)
		key = strings.Replace(string(newValue), "\"", "", -1)
	}

	return key
}

func FormatPrice(Price string) string {
	var result string
	tmp := strings.Split(Price, ".")
	if len(tmp) <= 1 {
		result = tmp[0] + ".00"
	} else {
		var i int
		total := 2 - utf8.RuneCountInString(tmp[1])
		if utf8.RuneCountInString(tmp[1]) > 2 {
			result = tmp[0] + "." + tmp[1][0:2]
		} else {
			result = tmp[0] + "." + tmp[1]
		}
		for i = 0; i < total; i++ {
			result = result + "0"
		}
	}
	return result
}

//反转字符串
func Reverse(s string) string {
	a := func(s string) *[]rune {
		var b []rune
		for _, k := range []rune(s) {
			defer func(v rune) {
				b = append(b, v)
			}(k)
		}
		return &b
	}(s)
	return string(*a)
}
