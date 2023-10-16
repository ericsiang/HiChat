package initialize

import (
	"log"

	"go.uber.org/zap"
)


func InitLogger(){
	 //初始化日志
	logger ,err := zap.NewProduction()
	if err != nil{
		log.Fatal("[Logger Init Error] ",err.Error())
	}

	//使用全局logger
	zap.ReplaceGlobals(logger)
}