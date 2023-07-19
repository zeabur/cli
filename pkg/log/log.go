// Package log contains the loggers for the CLI
package log

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

// NewDebugLevel returns a logger with debug level
func NewDebugLevel() *zap.SugaredLogger {
	return New(zapcore.DebugLevel)
}

// NewInfoLevel returns a logger with info level
func NewInfoLevel() *zap.SugaredLogger {
	return New(zapcore.InfoLevel)
}

// New returns a logger with the given level
func New(level zapcore.Level) *zap.SugaredLogger {
	conf := zap.NewDevelopmentConfig()
	conf.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	conf.EncoderConfig.EncodeTime = func(time.Time, zapcore.PrimitiveArrayEncoder) {}
	conf.Level = zap.NewAtomicLevelAt(level)

	logger, _ := conf.Build()

	return zap.New(logger.Core()).Sugar()
}

// NewForUT returns a logger with the given level and buffer for unit testing
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
