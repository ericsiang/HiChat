package initialize

import (
	"HiChat/global"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(){
	//注意：User和Password为MySQL数据库的管理员密码，Host和Port为数据库连接ip端口，DBname为要连接的数据库
	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", User,Password,Ip,Port,DBName)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", global.ServiceConfig.MysqlDB.User,global.ServiceConfig.MysqlDB.Password,global.ServiceConfig.MysqlDB.Host,global.ServiceConfig.MysqlDB.Port,global.ServiceConfig.MysqlDB.DBName)

	newLogger := logger.New(
		log.New(os.Stdout,"\r\n",log.LstdFlags), // io writer （日志输出的目标，前缀和日志包含的内容——译者注）

		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  true,        // 禁用彩色打印
		},
	)

	var err error

	//将获取到的连接赋值到global.DB
	global.DB ,err =gorm.Open(mysql.Open(dsn),&gorm.Config{
		Logger: newLogger, //打印sql日志
	})

	if err != nil{
		panic(err)
	}
}