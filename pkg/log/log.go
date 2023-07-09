package log

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

func NewDebugLevel() *zap.SugaredLogger {
	return New(zapcore.DebugLevel)
}

func NewInfoLevel() *zap.SugaredLogger {
	return New(zapcore.InfoLevel)
}

func New(level zapcore.Level) *zap.SugaredLogger {
	conf := zap.NewDevelopmentConfig()
	conf.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	conf.EncoderConfig.EncodeTime = func(time.Time, zapcore.PrimitiveArrayEncoder) {}
	conf.Level = zap.NewAtomicLevelAt(level)

	logger, _ := conf.Build()

	return zap.New(logger.Core()).Sugar()
}

func NewForUT(buffer *zaptest.Buffer, level zapcore.Level) *zap.SugaredLogger {
	conf := zap.NewDevelopmentEncoderConfig()
	conf.EncodeLevel = zapcore.CapitalLevelEncoder // note: without color
	conf.EncodeTime = func(time.Time, zapcore.PrimitiveArrayEncoder) {}

	logger := zap.NewExample().WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewCore(
			zapcore.NewConsoleEncoder(conf),
			zapcore.AddSync(buffer),
			level,
		)
	}))

	return zap.New(logger.Core()).Sugar()
}
