package redis

import (
	"fmt"
	"testing"
)

func Test_Set(t *testing.T) {
	//GetInstance("ppt").Set("kkk","vvv",3600)
	//GetInstance("ppt").MSetJson("k1","v1","k2","v2","k3","v3")
	//GetInstance("ppt").Expire(3600*3,"k1","k2","k3")
	//res := GetInstance("ppt").MGet("k1","k2","k3")
	//fmt.Println(res)
	//GetInstance("ppt").Get("k1")
	//GetInstance("ppt").Keys("k*")
	//GetInstance("ppt").HMSet("ddd", map[string]interface{}{"rrr":35})
	//d := []HMSMD{}
	//d = append(d,HMSMD{Key:"wwwww",Data: map[string]interface{}{"ww":123},Ttl:3600})
	//d = append(d,HMSMD{Key:"ddd",Data: map[string]interface{}{"ww":123},Ttl:3600})
	//d = append(d,HMSMD{Key:"ccc",Data: map[string]interface{}{"ww":123},Ttl:3600})
	//GetInstance("ppt").HMSetMulti(d)
	//res := GetInstance("ppt").HMGet("aaa","111","222")
	//res := GetInstance("ppt").HMGetMulti(map[string][]string{"aaa":[]string{"111","222"},"bbb":[]string{"333","444"}})
	//res := GetInstance("ppt").HGetAll("aaa")
	//res := GetInstance("ppt").HDel("aaa","New Key")
	//res := GetInstance("ppt").Delete("aaa","bbb")
	//res := GetInstance("ppt").Incr("inum")
	//res := GetInstance("ppt").IncrBy("inum",21)
	//res := GetInstance("ppt").DecrByMulti(map[string]int{"inum":2,"inum2":10})
	//res := GetInstance("ppt").IncrByMulti(map[string]int{"inum":2,"inum2":10})
	//res := GetInstance("ppt").ExistsMulti("hhh","inum")
	//res := GetInstance("ppt").HIncr("aaa","111",29)
	//res := GetInstance("ppt").HExists("aaa","111")
	//res := GetInstance("ppt").Zadd("aaa1",map[int]string{1:"num",9:"name",100:"age"})
	//res := GetInstance("ppt").Zrange("aaa1",0,2)
	//res := GetInstance("ppt").Zrem("aaa1", "num", "age")
	//fmt.Println
	//dkey := []map[string][]string{}
	//karr := []string{"ww"}
	//delkey := map[string][]string{"ddd": karr}
	//dkey = append(dkey, delkey)
	res := GetInstance("ppt").DumpHGetAll("这是个测试id")
	fmt.Println(res)

}
