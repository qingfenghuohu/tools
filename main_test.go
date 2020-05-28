package tools

import (
	"fmt"
	"github.com/qingfenghuohu/tools/redis"
	"testing"
)

func TestMtRand1_Post(t *testing.T) {
	res := redis.GetInstance("ppt").Zrange("a", 0, 4)
	//res := redis.GetInstance("ppt").Zrevrange("a",2,4)
	//data := map[int]string{6:"6a",7:"7a",8:"8a",9:"9a",5:"5a"}
	//res := redis.GetInstance("ppt").Zadd("a",data)
	fmt.Println(res)
}
func TestMtRand_Post(t *testing.T) {
	//res := redis.GetInstance("ppt").Exists("hhh")
	//res := redis.GetInstance("ppt").HDecr("aaa", "111", 3)
	res := redis.GetInstance("ppt").HExists("aaa", "112")
	fmt.Println(res)
	//data := []interface{}{}
	//data = append(data, "hhh")
	//data = append(data, "1")
	//data = append(data, "hhh")
	//data = append(data, "2")
	//redis.GetInstance("ppt").HMDel(data)
	//d1 := redis.HMSMD{"hhh", map[string]interface{}{"2":true},86400}
	//d4 := redis.HMSMD{"hhh", map[string]interface{}{"3":false},86400}
	//d2 := redis.HMSMD{"aaa", map[string]interface{}{"1":"iiiii"},86400}
	//d3 := redis.HMSMD{"hhh", map[string]interface{}{"1":123},86400}
	////
	//data := []redis.HMSMD{}
	//data = append(data, d1)
	//data = append(data, d2)
	//data = append(data, d3)
	//data = append(data, d4)
	//redis.GetInstance("ppt").HMSetMulti(data)

	//data := map[string][]string{
	//	"hhh": []string{"1", "2", "3"},
	//	"aaa": []string{"1", "2", "3"},
	//}
	//result := redis.GetInstance("ppt").HMGetMulti(data)
	//fmt.Println(result)
	//redis.GetInstance("ppt").DecrBy("aaa", 3)
}
