package initialize

import (
	"HiChat/global"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func InitConfig(){
	v := viper.New()
	configFile := "../HiChat/config-debug.yaml"

	//讀取配置文件
	v.SetConfigFile(configFile)

	//讀取配置信息
	if err := v.ReadInConfig();err !=nil{
		panic(err)
	}

	//將數據放入全局變量 global.ServiceConfig
	if err := v.Unmarshal(&global.ServiceConfig);err !=nil{
		panic(err)
	}

	zap.S().Info("配置信息 : " , global.ServiceConfig)
}