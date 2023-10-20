package initialize

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitCoreConfig() zapcore.Core {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",   //日誌時間的key
		LevelKey:       "level",  //日誌級別的key
		NameKey:        "logger", //日誌名的keyㄌ
		MessageKey:     "msg",    //日誌消息的key
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,     //日誌結尾分隔符 - 默認/n
		EncodeLevel:    zapcore.LowercaseLevelEncoder, //日志级别，默認小寫
		EncodeTime:     zapcore.ISO8601TimeEncoder,    // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 全路径编码器
	}
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	//依不同級別寫入不同文件
	//info level log
	infoLoggerWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./log/info/info.log", //日誌文件存放目錄，如果文件夾不存在會自動創建
		MaxSize:    1,                     //文件大小限制,單位MB
		MaxBackups: 5,                     //最大保留日誌文件數量
		MaxAge:     5,                     //日誌文件保留天數
		Compress:   false,                 //是否壓縮處理
	})


	infoCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(infoLoggerWriteSyncer, zapcore.AddSync(os.Stdout)), zapcore.InfoLevel)

	//error level log
	errorLoggerWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./log/error/error.log", //日誌文件存放目錄，如果文件夾不存在會自動創建
		MaxSize:    1,                       //文件大小限制,單位MB
		MaxBackups: 5,                       //最大保留日誌文件數量
		MaxAge:     5,                       //日誌文件保留天數
		Compress:   false,                   //是否壓縮處理
	})

	errorCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(errorLoggerWriteSyncer, zapcore.AddSync(os.Stdout)), zapcore.ErrorLevel)

	teeCore := zapcore.NewTee(infoCore, errorCore)

	return teeCore
}



func InitLogger() {
	//初始化zap日志
	teeCore := InitCoreConfig()
	logger := zap.New(teeCore, zap.AddCaller(),zap.Development())

	defer logger.Sync() // zap底层有缓冲。在任何情况下执行 defer logger.Sync() 是一个很好的习惯
	
	zap.ReplaceGlobals(logger) //使用全局logger(設定了在其他地方調用 zap.S() or zap.L() 才會生效)
}
