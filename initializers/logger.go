package initializers

import (
	"os"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type CustomLogLevel zapcore.Level

const (
	InfoLevel  CustomLogLevel = CustomLogLevel(zapcore.InfoLevel)
	WarnLevel  CustomLogLevel = CustomLogLevel(zapcore.WarnLevel)
	ErrorLevel CustomLogLevel = CustomLogLevel(zapcore.ErrorLevel)
)

func AddLogger() *zap.SugaredLogger {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	level := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel
	})

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(config.EncoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),
		level,
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return logger.Sugar() 
}

var Logger *zap.SugaredLogger

func init() {
	Logger = AddLogger()
}