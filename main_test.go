package tools

import (
	"fmt"
	"github.com/qingfenghuohu/tools/redis"
	"testing"
)

func TestMtRand_Post(t *testing.T) {
	//d1 := redis.HMSMD{"hhh", map[string]interface{}{"2":true},86400}
	//d4 := redis.HMSMD{"hhh", map[string]interface{}{"3":false},86400}
	//d2 := redis.HMSMD{"aaa", map[string]interface{}{"1":"iiiii"},86400}
	//d3 := redis.HMSMD{"hhh", map[string]interface{}{"1":123},86400}
	//
	//data := []redis.HMSMD{}
	//data = append(data, d1)
	//data = append(data, d2)
	//data = append(data, d3)
	//data = append(data, d4)
	//redis.GetInstance("ppt").HMSetMulti(data)

	data := map[string][]string{
		"hhh": []string{"1", "2", "3"},
		"aaa": []string{"1", "2", "3"},
	}
	result := redis.GetInstance("ppt").HMGetMulti(data)
	fmt.Println(result)
	//redis.GetInstance("ppt").DecrBy("aaa", 3)
}
