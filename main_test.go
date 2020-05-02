package tools

import (
	"github.com/qingfenghuohu/tools/redis"
	"testing"
)

func TestMtRand_Post(t *testing.T) {
	redis.GetInstance("ppt").IncrBy("aaa", 11)
	redis.GetInstance("ppt").DecrBy("aaa", 3)
}
