package tools

import "tools/redis"

func main() {
	redis.GetInstance("joe").MGet("1:i:joe.domain:state:")
}
