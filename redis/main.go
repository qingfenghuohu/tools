package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/json-iterator/go"
	"github.com/qingfenghuohu/config"
	"github.com/qingfenghuohu/tools/str"
	"net"
	"strconv"
	"sync"
	"time"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary
var GClient *redis.Client

var (
	DEFAULT = time.Duration(0)  // 过期时间 不设置
	FOREVER = time.Duration(-1) // 过期时间不设置
)

var Connect map[string]Cache

type Cache struct {
	Conn *redis.Client
}

var lock sync.Mutex

func GetInstance(dbName string) Cache {
	if len(Connect) == 0 {
		Connect = map[string]Cache{}
	}
	if _, ok := Connect[dbName]; !ok {
		lock.Lock()
		defer lock.Unlock()
		redisConfig := config.Data["cache"].(map[string]interface{})["redis"].(map[string]interface{})[dbName].(map[string]interface{})
		host := redisConfig["host"].(string) + ":" + redisConfig["port"].(string)
		db, _ := strconv.Atoi(redisConfig["db"].(string))
		pwd, _ := redisConfig["pwd"].(string)
		Connect[dbName] = Cache{Conn: NewRedisCache(db, host, pwd, 30)}
	}
	return Connect[dbName]
}

// 返回cache 对象, 在多个工具之间建立一个 中间初始化的时候使用
func NewRedisCache(db int, host string, pwd string, defaultExpiration time.Duration) *redis.Client {
	return redis.NewClient(&redis.Options{
		//连接信息
		Network:  "tcp", //网络类型，tcp or unix，默认tcp
		Addr:     host,  //主机名+冒号+端口，默认localhost:6379
		Password: pwd,   //密码
		DB:       db,    // redis数据库index
		//连接池容量及闲置连接数量
		PoolSize:     150, // 连接池最大socket连接数，默认为4倍CPU数， 4 * runtime.NumCPU
		MinIdleConns: 10,  //在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量；。
		//超时
		DialTimeout:  defaultExpiration, //连接建立超时时间，默认5秒。
		ReadTimeout:  3 * time.Second,   //读超时，默认3秒， -1表示取消读超时
		WriteTimeout: 3 * time.Second,   //写超时，默认等于读超时
		PoolTimeout:  4 * time.Second,   //当所有连接都处在繁忙状态时，客户端等待可用连接的最大等待时长，默认为读超时+1秒。
		//闲置连接检查包括IdleTimeout，MaxConnAge
		IdleCheckFrequency: 60 * time.Second, //闲置连接检查的周期，默认为1分钟，-1表示不做周期性检查，只在客户端获取连接时对闲置连接进行处理。
		IdleTimeout:        5 * time.Minute,  //闲置超时，默认5分钟，-1表示取消闲置超时检查
		MaxConnAge:         0 * time.Second,  //连接存活时长，从创建开始计时，超过指定时长则关闭连接，默认为0，即不关闭存活时长较长的连接
		//命令执行失败时的重试策略
		MaxRetries:      0,                      // 命令执行失败时，最多重试多少次，默认为0即不重试
		MinRetryBackoff: 8 * time.Millisecond,   //每次计算重试间隔时间的下限，默认8毫秒，-1表示取消间隔
		MaxRetryBackoff: 512 * time.Millisecond, //每次计算重试间隔时间的上限，默认512毫秒，-1表示取消间隔
		//可自定义连接函数
		Dialer: func() (net.Conn, error) {
			netDialer := &net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 5 * time.Minute,
			}
			return netDialer.Dial("tcp", "127.0.0.1:6379")
		},
		//钩子函数
		OnConnect: func(conn *redis.Conn) error { //仅当客户端执行命令时需要从连接池获取连接时，如果连接池需要新建连接时则会调用此钩子函数
			//fmt.Printf("conn=%v\n", conn)
			return nil
		},
	})
}

// string 类型 添加, v 可以是任意类型
func (c Cache) Set(name string, v interface{}, ttl int64) bool {
	res := c.Conn.Set(name, v, time.Duration(1000*1000*1000*ttl))
	result, err := res.Result()
	if err != nil {
		fmt.Println(err.Error())
	}
	if result == "OK" {
		return true
	} else {
		return false
	}
}

// string 类型 添加, v 可以是任意类型
func (c Cache) MSetJson(args ...interface{}) bool {
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
	res := c.Conn.MSet(params...)
	result, err := res.Result()
	if err != nil {
		fmt.Println(err.Error())
	}
	if result == "OK" {
		return true
	} else {
		return false
	}
}

// string 类型 添加, v 可以是任意类型
func (c Cache) MSet(params ...interface{}) bool {
	res := c.Conn.MSet(params...)
	result, err := res.Result()
	if err != nil {
		return false
	}
	if result == "OK" {
		return true
	} else {
		return false
	}
}

// 设置过期时间 （单位 秒）
func (c Cache) Expire(LifeTime int64, keys ...string) {
	// 设置key 的过期时间
	pipe := c.Conn.Pipeline()
	for _, v := range keys {
		pipe.Expire(v, time.Duration(LifeTime*1000*1000*1000))
	}
	res, err := pipe.Exec()
	for k, v := range res {
		fmt.Println("key", k)
		fmt.Println("val", v)
	}
	fmt.Println(res, err)
}

func (c Cache) MGet(args ...string) map[string]interface{} {
	result := make(map[string]interface{})
	params := []string{}
	for _, v := range args {
		params = append(params, v)
	}
	res := c.Conn.MGet(params...)
	res1, err := res.Result()
	if err != nil {
		fmt.Println(err.Error())
	}
	for k, v := range params {
		result[v] = res1[k]
	}
	return result
}

// 获取 字符串类型的值
func (c Cache) Get(name string) string {
	res := c.Conn.Get(name)
	result, err := res.Result()
	if err != nil {
		fmt.Println(err.Error())
	}
	return result
}

func (c Cache) Keys(name string) []string {
	var result []string
	res := c.Conn.Keys(name)
	result, err := res.Result()
	if err != nil {
		fmt.Println(err.Error())
	}
	return result
}

func (c Cache) HSet(name string, field string, value interface{}) bool {
	var result bool
	res := c.Conn.HSet(name, field, value)
	result, err := res.Result()
	if err != nil {
		fmt.Println(err.Error())
	}
	return result
}

func (c Cache) HGet(name string, field string) string {
	res := c.Conn.HGet(name, field)
	result, err := res.Result()
	if err != nil {
		fmt.Println(err.Error())
	}
	return result
}

func (c Cache) HMSet(name string, args map[string]interface{}) bool {
	res := c.Conn.HMSet(name, args)
	result, err := res.Result()
	if err != nil {
		fmt.Println(err.Error())
	}
	if result == "OK" {
		return true
	}
	return false
}

type HMSMD struct {
	Key  string
	Data map[string]interface{}
	Ttl  int
}

func (c Cache) HMSetMulti(data []HMSMD) {
	if len(data) == 0 {
		return
	}
	pipe := c.Conn.Pipeline()
	for _, val := range data {
		pipe.HMSet(val.Key, val.Data)
		pipe.Expire(val.Key, time.Duration(val.Ttl*1000*1000*1000))
	}
	_, err := pipe.Exec()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func (c Cache) HMGet(name string, args ...string) map[string]string {
	result := map[string]string{}
	res := c.Conn.HMGet(name, args...)
	res1, err := res.Result()
	if err != nil {
		fmt.Println(err.Error())
	}
	for k, v := range args {
		result[v] = str.Obj2Str(res1[k])
	}
	return result
}

func (c Cache) HMGetMulti(data map[string][]string) map[string]map[string]string {
	result := map[string]map[string]string{}
	pipe := c.Conn.Pipeline()
	for name, val := range data {
		pipe.HMGet(name, val...)
	}
	res, err := pipe.Exec()
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, value := range res {
		for k, v := range value.Args() {
			key := value.Args()[1].(string)
			if _, ok := result[key]; !ok {
				result[key] = map[string]string{}
			}
			if value.Args()[0] == "hmget" && k != 0 && k != 1 {
				res1, err := value.(*redis.SliceCmd).Result()
				if err != nil {
					fmt.Println(err.Error())
				}
				for kk, vv := range res1 {
					if k-2 == kk {
						result[key][v.(string)] = str.Obj2Str(vv)
					}
				}
			}
		}
	}
	return result
}
func (c Cache) HGetAll(name string) map[string]string {
	res := c.Conn.HGetAll(name)
	result, err := res.Result()
	if err != nil {
		fmt.Println(err.Error())
	}
	return result
}
func (c Cache) HDel(name string, field ...string) bool {
	res := c.Conn.HDel(name, field...)
	result, err := res.Result()
	if err != nil {
		fmt.Println(err.Error())
	}
	if result > 0 {
		return true
	}
	return false
}
func (c Cache) HLen(name string) int {
	res := c.Conn.HLen(name)
	result, err := res.Result()
	if err != nil {
		fmt.Println(err.Error())
	}
	return int(result)
}

// 删除指定的键
func (c Cache) Delete(keys ...string) bool {
	res := c.Conn.Del(keys...)
	result, err := res.Result()
	if err != nil {
		fmt.Println(err.Error())
	}
	if result > 0 {
		return true
	}
	return false
}
func (c Cache) Incr(key string) int {
	res := c.Conn.Incr(key)
	result, err := res.Result()
	if err != nil {
		fmt.Println(err.Error())
	}
	return int(result)
}
func (c Cache) Decr(key string) int {
	res := c.Conn.Decr(key)
	result, err := res.Result()
	if err != nil {
		fmt.Println(err.Error())
	}
	return int(result)
}

func (c Cache) DecrBy(name string, num int) int {
	res := c.Conn.DecrBy(name, int64(num))
	result, err := res.Result()
	if err != nil {
		fmt.Println(err.Error())
	}
	return int(result)
}

func (c Cache) IncrBy(name string, num int) int {
	res := c.Conn.IncrBy(name, int64(num))
	result, err := res.Result()
	if err != nil {
		fmt.Println(err.Error())
	}
	return int(result)
}

func (c Cache) DecrByMulti(data map[string]int) map[string]int {
	result := map[string]int{}
	pipe := c.Conn.Pipeline()
	for k, v := range data {
		pipe.DecrBy(k, int64(v))
	}
	res, err := pipe.Exec()
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, v := range res {
		key := v.Args()[1].(string)
		res1, err := v.(*redis.IntCmd).Result()
		if err != nil {
			fmt.Println(err.Error())
		}
		result[key] = int(res1)
	}
	return result
}

func (c Cache) IncrByMulti(data map[string]int) map[string]int {
	result := map[string]int{}
	pipe := c.Conn.Pipeline()
	for k, v := range data {
		pipe.IncrBy(k, int64(v))
	}
	res, err := pipe.Exec()
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, v := range res {
		key := v.Args()[1].(string)
		res1, err := v.(*redis.IntCmd).Result()
		if err != nil {
			fmt.Println(err.Error())
		}
		result[key] = int(res1)
	}
	return result
}

func (c Cache) Exists(key string) bool {
	res := c.Conn.Exists(key)
	result, err := res.Result()
	if err != nil {
		fmt.Println(err.Error())
	}
	if result > 0 {
		return true
	}
	return false
}

func (c Cache) ExistsMulti(keys ...string) map[string]bool {
	result := map[string]bool{}
	pipe := c.Conn.Pipeline()
	for _, v := range keys {
		pipe.Exists(v)
	}
	res, err := pipe.Exec()
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, v := range res {
		key := v.Args()[1].(string)
		res1, err := v.(*redis.IntCmd).Result()
		if err != nil {
			fmt.Println(err.Error())
		}
		if res1 > 0 {
			result[key] = true
		} else {
			result[key] = false
		}
	}
	return result
}
func (c Cache) HIncr(key, field string, val int) int {
	var result int64
	res := c.Conn.HIncrBy(key, field, int64(val))
	result, err := res.Result()
	if err != nil {
		fmt.Println(err.Error())
	}
	return int(result)
}
func (c Cache) HDecr(key, field string, val int) int {
	var result int64
	val = 0 - val
	res := c.Conn.HIncrBy(key, field, int64(val))
	result, err := res.Result()
	if err != nil {
		fmt.Println(err.Error())
	}
	return int(result)
}
func (c Cache) HExists(key, field string) bool {
	var result bool
	res := c.Conn.HExists(key, field)
	result, err := res.Result()
	if err != nil {
		fmt.Println(err.Error())
	}
	return result
}
func (c Cache) Zadd(key string, data map[int]string) bool {
	var result int64
	params := []redis.Z{}
	for k, v := range data {
		params = append(params, redis.Z{Score: float64(k), Member: v})
	}
	res := c.Conn.ZAdd(key, params...)
	result, err := res.Result()
	if err != nil {
		fmt.Println(err.Error())
	}
	if result > 0 {
		return true
	}
	return false
}
func (c Cache) Zrange(key string, start, end int) []string {
	res := c.Conn.ZRange(key, int64(start), int64(end))
	result, err := res.Result()
	if err != nil {
		fmt.Println(err.Error())
	}
	return result
}
func (c Cache) Zrevrange(key string, start, end int) []string {
	res := c.Conn.ZRevRange(key, int64(start), int64(end))
	result, err := res.Result()
	if err != nil {
		fmt.Println(err.Error())
	}
	return result
}
func (c Cache) Zrem(key string, member ...string) bool {
	res := c.Conn.ZRem(key, member)
	result, err := res.Result()
	if err != nil {
		fmt.Println(err.Error())
	}
	if int(result) >= len(member) {
		return true
	}
	return false
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
