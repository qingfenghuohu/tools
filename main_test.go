package tools

import (
	"fmt"
	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/qingfenghuohu/tools/curl"
	"github.com/qingfenghuohu/tools/redis"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func TestMtRand1_Post(t *testing.T) {
	s := "https://www.jjppt.com/d/7862"
	u, _ := url.Parse(s)
	var Type int
	var Val string
	var matched [][]string
	if strings.Index(u.Path, "/beijing/detail/") == 0 {
		Type = 2
		re := regexp.MustCompile(`\/beijing\/detail\/([\d]+)`)
		matched = re.FindAllStringSubmatch(u.Path, -1)

	}
	if strings.Index(u.Path, "/jianli/detail/") == 0 {
		Type = 3
		re := regexp.MustCompile(`\/jianli\/detail\/([\d]+)`)
		matched = re.FindAllStringSubmatch(u.Path, -1)
	}
	if strings.Index(u.Path, "/d/") == 0 {
		Type = 1
		re := regexp.MustCompile(`\/d\/([\d]+)`)
		matched = re.FindAllStringSubmatch(u.Path, -1)
	}
	if len(matched) > 0 {
		if len(matched[0]) == 2 {
			Val = matched[0][1]
		}
	}
	ReqUrl := fmt.Sprintf("https://www.jjppt.com/vip/download?id=%s&type=%d", Val, Type)
	Curl := curl.New{}
	CookieData := map[string]interface{}{
		"AGL_USER_ID":        "f1462d05-45f8-4fa4-abbf-74c1eb3c6c86",
		"UM_distinctid":      "16f837a4459657-015a3d1f918bae-37647e05-13c680-16f837a445a3c6",
		"CNZZDATA1272242338": "52898249-1578455646-https%253A%252F%252Fwww.jjppt.com%252F%7C1578488169",
		"CNZZDATA1278254644": "1381914079-1578458175-%7C1578622520",
		"advanced-frontend":  "6h2grcak5qq5shninr54o3go05",
		"_identity-frontend": "516d655404e7ce0bb9fdabca241faefe364d4ff92a58b35f06f2744da21f00b9a%3A2%3A%7Bi%3A0%3Bs%3A18%3A%22_identity-frontend%22%3Bi%3A1%3Bs%3A50%3A%22%5B1364574%2C%22ITaaiC2jP_to0DYbViG_alnyHpJthLar%22%2C86400%5D%22%3B%7D",
	}
	random := browser.Random()
	random = random + "  " + strconv.Itoa(MtRand(100000, 999999))
	HeaderData := map[string]interface{}{
		"User-Agent":      random,
		"X-Forwarded-For": GenIpaddr(),
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3",
		"Accept-Encoding": "gzip, deflate, br",
		"Accept-Language": "zh-CN,zh;q=0.9",
		"Cache-Control":   "no-cache",
	}
	result := Curl.Cookie(CookieData).Header(HeaderData).Get(ReqUrl)
	fmt.Println(result)

	//res := redis.GetInstance("ppt").Zrange("a", 0, 4)
	//res := redis.GetInstance("ppt").Zrevrange("a",2,4)
	//data := map[int]string{6:"6a",7:"7a",8:"8a",9:"9a",5:"5a"}
	//res := redis.GetInstance("ppt").Zadd("a",data)
	//res := redis.GetInstance("ppt").Zrem("a", "9a", "8a")
	//fmt.Println(res)
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
