package tools

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/json-iterator/go"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

func Md5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}
func MtRand(min int, max int) int {
	rand.Seed(time.Now().Unix())
	return min + rand.Intn(max-min)
}
func JsonToMap(content interface{}) map[string]interface{} {
	var name map[string]interface{}
	if marshalContent, err := json.Marshal(content); err != nil {
		fmt.Println(err)
	} else {
		d := json.NewDecoder(bytes.NewReader(marshalContent))
		d.UseNumber() // 设置将float64转为一个number
		if err := d.Decode(&name); err != nil {
			fmt.Println(err)
		} else {
			for k, v := range name {
				name[k] = v
			}
		}
	}
	return name
}

//阻塞式的执行外部shell命令的函数,等待执行完毕并返回标准输出
func ExecShell(s string) (string, error) {
	//函数返回一个*Cmd，用于使用给出的参数执行name指定的程序
	cmd := exec.Command("/bin/bash", "-c", s)

	//读取io.Writer类型的cmd.Stdout，再通过bytes.Buffer(缓冲byte类型的缓冲器)将byte类型转化为string类型(out.String():这是bytes类型提供的接口)
	var out bytes.Buffer
	cmd.Stdout = &out

	//Run执行c包含的命令，并阻塞直到完成。  这里stdout被取出，cmd.Wait()无法正确获取stdin,stdout,stderr，则阻塞在那了
	err := cmd.Run()

	return out.String(), err
}

func CreateFile(FileName string, Content string) {
	f, err := os.Create(FileName)
	defer f.Close()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		_, err = f.Write([]byte(Content))
		CheckError(err)
	}
}
func SaveFile(FileName string, Content []byte) {
	f, err := os.Create(FileName)
	defer func() {
		if err := f.Close(); err != nil {
			// log etc
		}
	}()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		_, err = f.Write(Content)
		CheckError(err)
	}
}

func CheckError(err error) {

}

func RemoveDuplicateElement(addrs []string) []string {
	result := make([]string, 0, len(addrs))
	temp := map[string]struct{}{}
	for _, item := range addrs {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func UrlEncode(str string) string {
	return url.QueryEscape(str)
}

func KsortPostForm(params map[string][]string) string {
	var dataParams string
	//ksort
	var keys []string
	for k, _ := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	//拼接
	for _, k := range keys {
		dataParams = dataParams + k + params[k][0]
	}

	return dataParams
}

func EnCode(id int) string {
	return strconv.Itoa(MtRand(11111, 99999)) + strconv.Itoa((id+64)*16) + strconv.Itoa(MtRand(111, 999))
}

func DeCode(code string) int {
	s := string([]byte(code)[5 : len(code)-3])
	number, _ := strconv.Atoi(s)
	return (number / 16) - 64
}

func RandSeq(n int) string {
	letters := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func Interface2MapStrStr(m interface{}) map[string]map[string]string {
	var result map[string]map[string]string
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	data, _ := json.Marshal(m)
	reader := strings.NewReader(string(data))
	decoder := json.NewDecoder(reader)
	decoder.Decode(&result)
	return result
}

func Struct2Map(m interface{}) map[string]string {
	var result map[string]string
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	data, _ := json.Marshal(m)
	fmt.Println(string(data))
	reader := strings.NewReader(string(data))
	decoder := json.NewDecoder(reader)
	decoder.Decode(&result)
	fmt.Println(result)
	return result
}

func Interface2MapSliceStr(m interface{}) map[string][]map[string]string {
	var result map[string][]map[string]string
	data, _ := json.Marshal(m)
	json.Unmarshal(data, &result)
	return result
}

func CacheDataFormat(m interface{}) map[string]map[string]interface{} {
	var result map[string]map[string]interface{}
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	data, _ := json.Marshal(m)
	fmt.Println(string(data))
	reader := strings.NewReader(string(data))
	decoder := json.NewDecoder(reader)
	decoder.Decode(&result)
	fmt.Println(result)
	return result
}

func HttpClientPost(url string, data string) string {
	resp, err := http.Post(url,
		"application/x-www-form-urlencoded",
		strings.NewReader(data))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	return string(body)
}

func GetHost(context *gin.Context) (string, string) {
	scheme := "https"
	host := context.Request.Host
	if strings.Index(host, "127.0.0.1") != -1 {
		host = "api.joess.online"
		scheme = "https"
	}
	if strings.Index(host, "localhost") != -1 {
		host = "localhost"
		scheme = "http"
	}
	return scheme, host
}

func StrToTime(str string) int64 {
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation("2006-01-02 15:04:05", str, loc)
	result := theTime.Unix()
	return result
}
func GenIpaddr() string {
	rand.Seed(time.Now().Unix())
	ip := fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
	return ip
}

func GetFieldName(structName interface{}) []string {
	t := reflect.TypeOf(structName)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		log.Println("Check type error not Struct")
		return nil
	}
	fieldNum := t.NumField()
	result := make([]string, 0, fieldNum)
	for i := 0; i < fieldNum; i++ {
		result = append(result, t.Field(i).Name)
	}
	return result
}

//func main() {
//	fmt.Println(redis.GetInstance("joe").MGet("1:i:joe.domain:state:"))
//}
