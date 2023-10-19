package main

import (
	"HiChat/global"
	"HiChat/initialize"
	"HiChat/router"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Pong(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
        "name":   "測試接口",
    })
}

func main() {
	//初始化配置
	initialize.InitConfig()
	//初始化日志
	initialize.InitLogger()
	//初始化数据库
	initialize.InitDB()
	initialize.InitRedis()

	
	router := router.Router()
	router.Run(fmt.Sprintf(":%d", global.ServiceConfig.Port))
}
