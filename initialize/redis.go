package initialize

import (
	"HiChat/global"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func InitRedis() {
	opt := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", global.ServiceConfig.Redis.Host, global.ServiceConfig.Redis.Port), // redis地址
		Password: "",  // redis密码，没有则留空
		DB:       10,
	}
	global.RedisDB = redis.NewClient(opt)
}
