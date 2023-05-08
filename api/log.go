package api

import (
	"os"
	"path/filepath"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level = zapcore.Level

// global zap logger
var globalLog *zap.SugaredLogger

func InitLog(dir string, level zapcore.Level) {

	if globalLog == nil {

		core := zapcore.NewCore(getEncoder(false), getWriter(dir, level), level)

		logger := zap.New(core, zap.AddCaller())

		globalLog = logger.Sugar()

	}

}

// FlushLog 刷新日志,这里没有校验 globalLog 是否是 nil
func FlushLog() {
	if globalLog != nil {
		globalLog.Sync()
	}

}

// D 刷新日志,这里没有校验 globalLog 是否是 nil
func LogD(args ...interface{}) {
	globalLog.Debug(args...)
}

// I 刷新日志,这里没有校验 globalLog 是否是 nil
func LogI(args ...interface{}) {
	globalLog.Info(args...)
}

// E 刷新日志,这里没有校验 globalLog 是否是 nil
func LogE(args ...interface{}) {
	globalLog.Error(args...)
}

// W 刷新日志,这里没有校验 globalLog 是否是 nil
func LogW(args ...interface{}) {
	globalLog.Warn(args...)
}

// F 刷新日志,这里没有校验 globalLog 是否是 nil
func LogF(args ...interface{}) {
	globalLog.Fatal(args...)
}

func getEncoder(json bool) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	if json {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getWriter(logDir string, lv Level) zapcore.WriteSyncer {

	if logDir == `stdout` {
		return zapcore.Lock(os.Stdout)
	}

	fileName := filepath.Join(logDir, "llmchain."+lv.CapitalString()+".log")

	lumberJackLogger := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	}
	return zapcore.AddSync(lumberJackLogger)
}
