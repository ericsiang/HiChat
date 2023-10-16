package main

import (
	"HiChat/initialize"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Pong(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
        "name":   "測試接口",
    })
}

func main() {
	//初始化日志
	initialize.InitLogger()
	//初始化数据库
	initialize.InitDB("127.0.0.1","root","123456","hiChat",3306)

	r:= gin.Default()
	r.GET("/ping", Pong) //测试接口
	r.Run(":8083")
}
