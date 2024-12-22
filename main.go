package main

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Logger *zap.Logger
)

func init() {
	Logger = LoadLog()
}

func LoadLog() *zap.Logger {
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(os.Stdout),
		zap.DebugLevel,
	)
	Logger = zap.New(core)
	defer Logger.Sync()
	return Logger
}

func WarpLogError(err error) zap.Field {
	return zap.Error(err)
}

func LogError(msg string, err error) {
	Logger.Error(msg, WarpLogError(err))
}

func main() {

}
