package log

import (
	"os"

	lumberjack "gopkg.in/natefinch/lumberjack.v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"usersvr/config"
)

var log *zap.Logger
var sugar *zap.SugaredLogger

func InitLog() {
	var coreArr []zapcore.Core

	// 获取编码器
	encoderConfig := zap.NewProductionEncoderConfig()            // NewJSONEncoder()输出json格式，NewConsoleEncoder()输出普通文本格式
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder        // 指定时间格式
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // 按级别显示不同颜色，不需要的话取值zapcore.CapitalLevelEncoder就可以了
	// encoderConfig.EncodeCaller = zapcore.FullCallerEncoder        //显示完整文件路径
	// TODO:此处可换成输出json格式
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	// 日志级别
	highPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { // error级别
		return lev >= zap.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { // info和debug级别,debug级别是最低的
		if config.GetGlobalConfig().LogConfig.Level == "debug" {
			return lev < zap.ErrorLevel && lev >= zap.DebugLevel
		} else {
			return lev < zap.ErrorLevel && lev >= zap.InfoLevel
		}
	})
	// info文件writeSyncer
	logConfig := config.GetGlobalConfig().LogConfig
	infoFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logConfig.LogPath + "info_" + logConfig.FileName, // 日志文件存放目录，
		MaxSize:    logConfig.MaxSize,                                // 文件大小限制,单位MB
		MaxBackups: logConfig.MaxBackups,                             // 最大保留日志文件数量
		MaxAge:     logConfig.MaxAge,                                 // 日志文件保留天数
		Compress:   false,                                            // 是否压缩处理
	})
	infoFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(infoFileWriteSyncer, zapcore.AddSync(os.Stdout)), lowPriority)
	// error文件writeSyncer
	errorFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logConfig.LogPath + "error_" + logConfig.FileName, // 日志文件存放目录
		MaxSize:    logConfig.MaxSize,                                 // 文件大小限制,单位MB
		MaxBackups: logConfig.MaxBackups,                              // 最大保留日志文件数量
		MaxAge:     logConfig.MaxAge,                                  // 日志文件保留天数
		Compress:   false,                                             // 是否压缩处理
	})
	errorFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(errorFileWriteSyncer, zapcore.AddSync(os.Stdout)), highPriority)

	coreArr = append(coreArr, infoFileCore)
	coreArr = append(coreArr, errorFileCore)
	log = zap.New(zapcore.NewTee(coreArr...), zap.AddCaller(), zap.AddCallerSkip(1)) // zap.AddCaller()为显示文件名和行号，可省略
	sugar = log.Sugar()
}

func Errorf(s string, v ...interface{}) {
	sugar.Errorf(s, v...)
}
