package curl

import (
	"fmt"
	"testing"
)

func TestNew_Get(t *testing.T) {
	Curl := New{}
	d := map[string]interface{}{
		"test": "sign",
		"sn":   "3333333333",
	}
	result := Curl.Data(d).Post("http://localhost:8002/test")
	fmt.Println(result)
}
