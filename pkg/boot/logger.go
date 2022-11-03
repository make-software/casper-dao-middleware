package boot

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger create new logger with json encoder with os.Stdout WriteSyncer
func NewLogger(logLevel zapcore.Level) {
	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionConfig().EncoderConfig),
		os.Stdout,
		logLevel,
	))
	zap.ReplaceGlobals(logger)
}
