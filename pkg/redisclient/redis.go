package redisclient

import (
	"github.com/go-redis/redis/extra/redisotel/v8"
	"github.com/go-redis/redis/v8"
)

var Rdb *redis.ClusterClient

// InitRedis 初始化Redis客户端
func InitRedis() {
	Rdb = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    []string{"******:6380"},
		Password: "******", // no password set
	})

	//添加tracing的Hook
	Rdb.AddHook(redisotel.NewTracingHook())
}
