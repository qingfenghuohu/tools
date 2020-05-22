package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/json-iterator/go"
	"github.com/qingfenghuohu/config"
	"github.com/qingfenghuohu/tools/str"
	"strconv"
	"sync"
	"time"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

var (
	DEFAULT = time.Duration(0)  // 过期时间 不设置
	FOREVER = time.Duration(-1) // 过期时间不设置
)

func init() {
	Connect = make(map[string]*Cache)
}

type Cache struct {
	pool              *redis.Pool
	defaultExpiration time.Duration
}

var Connect map[string]*Cache
var lock sync.Mutex

func GetInstance(dbName string) *Cache {
	if _, ok := Connect[dbName]; ok {
	} else {
		lock.Lock()
		defer lock.Unlock()
		redisConfig := config.Data["cache"].(map[string]interface{})["redis"].(map[string]interface{})[dbName].(map[string]interface{})
		host := redisConfig["host"].(string) + ":" + redisConfig["port"].(string)
		db, _ := strconv.Atoi(redisConfig["db"].(string))
		pwd, _ := redisConfig["pwd"].(string)
		//Connect[dbName] = conn(db, host)
		Connect[dbName] = NewRedisCache(db, host, pwd, 30)

	}
	return Connect[dbName]
}

// 返回cache 对象, 在多个工具之间建立一个 中间初始化的时候使用
func NewRedisCache(db int, host string, pwd string, defaultExpiration time.Duration) *Cache {
	pool := &redis.Pool{
		MaxActive:   100,                              //  最大连接数，即最多的tcp连接数，一般建议往大的配置，但不要超过操作系统文件句柄个数（centos下可以ulimit -n查看）
		MaxIdle:     10,                               // 最大空闲连接数，即会有这么多个连接提前等待着，但过了超时时间也会关闭。
		IdleTimeout: time.Duration(100) * time.Second, // 空闲连接超时时间，但应该设置比redis服务器超时时间短。否则服务端超时了，客户端保持着连接也没用
		Wait:        true,                             // 当超过最大连接数 是报错还是等待， true 等待 false 报错
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", host, redis.DialDatabase(db), redis.DialPassword(pwd))
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
	}
	return &Cache{pool: pool, defaultExpiration: defaultExpiration}
}

// string 类型 添加, v 可以是任意类型
func (c Cache) Set(name string, v interface{}) bool {
	conn := c.pool.Get()
	defer conn.Close()
	_, err := conn.Do("SET", name, v)
	if err == nil {
		return true
	} else {
		return false
	}
}

// string 类型 添加, v 可以是任意类型
func (c Cache) MSetJson(args ...interface{}) bool {
	conn := c.pool.Get()
	defer conn.Close()
	params := make([]interface{}, len(args))
	k := 1
	for i, v := range args {
		if k%2 == 0 {
			tmp, _ := Serialization(v)
			params[i] = string(tmp)
		} else {
			params[i] = v
		}
		k++
	}
	res1, err := redis.String(conn.Do("MSET", params...))
	if err != nil {
		return false
	}
	if res1 == "OK" {
		return true
	} else {
		return false
	}
}

// string 类型 添加, v 可以是任意类型
func (c Cache) MSet(params ...interface{}) bool {
	conn := c.pool.Get()
	defer conn.Close()
	res1, err := redis.String(conn.Do("MSET", params...))
	if err != nil {
		return false
	}
	if res1 == "OK" {
		return true
	} else {
		return false
	}
}

// 设置过期时间 （单位 秒）
func (c Cache) Expire(newSecondsLifeTime int64, keys ...string) {
	// 设置key 的过期时间
	conn := c.pool.Get()
	defer conn.Close()
	for _, v := range keys {
		conn.Send("EXPIRE", v, newSecondsLifeTime)
	}
	conn.Flush()
}

func (c Cache) MGet(args ...string) map[string]interface{} {
	conn := c.pool.Get()
	defer conn.Close()
	params := []interface{}{}
	for _, v := range args {
		params = append(params, v)
	}
	res1, err := conn.Do("MGET", params...)
	result := make(map[string]interface{})
	if res1 != nil {
		for i, v := range res1.([]interface{}) {
			if v != nil {
				tmp, _ := redis.String(v, err)
				result[args[i]] = tmp
			}
		}
	}
	return result
}

// 获取 字符串类型的值
func (c Cache) Get(name string) string {
	conn := c.pool.Get()
	defer conn.Close()
	res, _ := redis.String(conn.Do("Get", name))
	return res
}

func (c Cache) Keys(name string) []string {
	var result []string
	conn := c.pool.Get()
	defer conn.Close()
	res, _ := redis.ByteSlices(conn.Do("keys", name))
	for _, v := range res {
		result = append(result, string(v))
	}
	return result
}
func (c Cache) HSet(name interface{}, field interface{}, value interface{}) bool {
	conn := c.pool.Get()
	defer conn.Close()
	_, err := conn.Do("hset", name, field, value)
	if err != nil {
		fmt.Println("hmset error", err.Error())
		return false
	}
	return true
}
func (c Cache) HGet(name interface{}, field interface{}) string {
	result := ""
	conn := c.pool.Get()
	defer conn.Close()
	res, err := conn.Do("hget", name, field)
	if err != nil {
		fmt.Println("hmget failed", err.Error())
	} else {
		if res != nil {
			result = string(res.([]byte))
		}
	}
	return result
}
func (c Cache) HMSet(name interface{}, args ...interface{}) bool {
	var params []interface{}
	conn := c.pool.Get()
	defer conn.Close()
	params = append(params, name)
	params = append(params, args...)
	_, err := conn.Do("hmset", params...)
	if err != nil {
		fmt.Println("hmset error", err.Error())
		return false
	}
	return true
}

type HMSMD struct {
	Key  string
	Data map[string]interface{}
	Ttl  int
}

func (c Cache) HMSetMulti(data []HMSMD) {
	conn := c.pool.Get()
	defer conn.Close()
	if len(data) == 0 {
		return
	}
	ttl := map[string]int{}
	for _, val := range data {
		for k, v := range val.Data {
			conn.Send("hmset", val.Key, k, str.Obj2Str(v))
		}
		ttl[val.Key] = val.Ttl
	}
	for k, v := range ttl {
		conn.Send("EXPIRE", k, v)
	}
	conn.Flush()
}
func (c Cache) HMGet(name interface{}, args ...interface{}) map[string]string {
	result := map[string]string{}
	var params []interface{}
	conn := c.pool.Get()
	defer conn.Close()
	params = append(params, name)
	params = append(params, args...)
	res, err := conn.Do("hmget", params...)
	if err != nil {
		fmt.Println("hmget failed", err.Error())
	} else {
		for i, v := range res.([]interface{}) {
			result[args[i].(string)] = string(v.([]byte))
		}
	}
	return result
}
func (c Cache) HMGetMulti(data map[string][]string) map[string]map[string]string {
	result := map[string]map[string]string{}
	conn := c.pool.Get()
	defer conn.Close()
	if len(data) == 0 {
		return result
	}
	for key, val := range data {
		for _, v := range val {
			conn.Send("hmget", key, v)
		}
	}
	conn.Flush()
	for key, val := range data {
		if len(result[key]) == 0 {
			result[key] = map[string]string{}
		}
		for _, v := range val {
			res, _ := redis.Values(conn.Receive())
			result[key][v] = ""
			if res[0] != nil {
				result[key][v] = string(res[0].([]byte))
			}
		}
	}
	return result
}
func TypeOf(v interface{}) string {
	return fmt.Sprintf("%T", v)
}
func (c Cache) HGetAll(name string) map[string]string {
	conn := c.pool.Get()
	defer conn.Close()

	res, err := redis.StringMap(conn.Do("HGETALL", name))
	if err != nil {
		fmt.Println("hmget failed", err.Error())
	}
	return res
}
func (c Cache) HDel(name string, field string) bool {
	conn := c.pool.Get()
	defer conn.Close()
	_, err := conn.Do("hdel", name, field)
	if err != nil {
		fmt.Println("hdel failed", err.Error())
		return false
	}
	return true
}
func (c Cache) HMDel(data []interface{}) bool {
	conn := c.pool.Get()
	defer conn.Close()
	if len(data) == 0 {
		return false
	}
	_, err := conn.Do("hdel", data...)
	if err == nil {
		return true
	}
	return false
}

// 删除指定的键
func (c Cache) Delete(keys ...interface{}) (bool, error) {
	conn := c.pool.Get()
	defer conn.Close()
	v, err := redis.Bool(conn.Do("DEL", keys...))
	return v, err
}
func (c Cache) Incr(name string) int {
	conn := c.pool.Get()
	defer conn.Close()
	res, _ := redis.Int(conn.Do("INCR", name))
	return res
}
func (c Cache) Decr(name string) int {
	conn := c.pool.Get()
	defer conn.Close()
	res, _ := redis.Int(conn.Do("DECR", name))
	return res
}

func (c Cache) DecrBy(name string, num int) int {
	conn := c.pool.Get()
	defer conn.Close()
	res, _ := redis.Int(conn.Do("DECRBY", name, num))
	return res
}

func (c Cache) IncrBy(name string, num int) int {
	conn := c.pool.Get()
	defer conn.Close()
	res, _ := redis.Int(conn.Do("INCRBY", name, num))
	return res
}

func (c Cache) DecrByMulti(data map[string]int) map[string]bool {
	var result map[string]bool
	conn := c.pool.Get()
	defer conn.Close()
	for k, v := range data {
		conn.Send("DECRBY", k, v)
	}
	conn.Flush()
	for k, _ := range data {
		result[k], _ = redis.Bool(conn.Receive())
	}
	return result
}

func (c Cache) IncrByMulti(data map[string]int) map[string]bool {
	var result map[string]bool
	conn := c.pool.Get()
	defer conn.Close()
	for k, v := range data {
		conn.Send("DECINCRBYRBY", k, v)
	}
	conn.Flush()
	for k, _ := range data {
		result[k], _ = redis.Bool(conn.Receive())
	}
	return result
}
func (c Cache) Exists(key string) bool {
	var result bool
	conn := c.pool.Get()
	res, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}
	result = res
	return result
}
func (c Cache) HIncr(key, field string, val int) bool {
	var result bool
	conn := c.pool.Get()
	res, err := redis.Bool(conn.Do("HINCRBY", key, field, val))
	if err != nil {
		return false
	}
	result = res
	return result
}
func (c Cache) HDecr(key, field string, val int) bool {
	var result bool
	conn := c.pool.Get()
	res, err := redis.Bool(conn.Do("HDECRBY", key, field, val))
	if err != nil {
		return false
	}
	result = res
	return result
}
func Deserialization(data []byte, i *interface{}) (interface{}, error) {
	result := new(interface{})
	err := json.Unmarshal(data, result)
	return result, err
}
func Serialization(v interface{}) ([]byte, error) {
	result, err := json.Marshal(v)
	return result, err
}

// Convert json string to map
func JsonToMap(jsonStr string) (map[string]string, error) {
	m := make(map[string]string)
	err := json.Unmarshal([]byte(jsonStr), &m)
	if err != nil {
		fmt.Printf("Unmarshal with error: %+v\n", err)
		return nil, err
	}

	for k, v := range m {
		fmt.Printf("%v: %v\n", k, v)
	}

	return m, nil
}
