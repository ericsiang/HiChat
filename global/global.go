package global

import (
	"HiChat/config"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	DB      *gorm.DB
	RedisDB *redis.Client
	ServiceConfig *config.ServiceConfig
)
