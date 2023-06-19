package lib

import (
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// global config
var config *viper.Viper

// global logger
var logger *zap.SugaredLogger

func Init(confFilePath string) (*viper.Viper, *zap.SugaredLogger) {
	// config
	config = viper.New()
	config.SetConfigFile(confFilePath)
	err := config.ReadInConfig()
	if err != nil {
		panic("read config file fail: " + err.Error())
	}
	// logger
	var core zapcore.Core
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./logs/app.log",
		MaxSize:    100, // megabytes
		MaxBackups: 3,
		MaxAge:     7, // days
	})
	if config.GetString("app.env") == "prod" {
		core = zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), w, zap.WarnLevel)
	} else {
		ec := zap.NewDevelopmentEncoderConfig()
		core = zapcore.NewTee(zapcore.NewCore(
			zapcore.NewJSONEncoder(ec),
			w,
			zap.DebugLevel,
		), zapcore.NewCore(
			zapcore.NewConsoleEncoder(ec),
			zapcore.AddSync(os.Stderr),
			zap.DebugLevel,
		))
	}
	logger = zap.New(core).Sugar()
	zap.RedirectStdLog(logger.Desugar())
	return config, logger
}
