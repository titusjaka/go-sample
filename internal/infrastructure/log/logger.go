package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger used for default-way logging in the project
type Logger interface {
	Info(message string, fields ...FieldAny)
	Error(message string, fields ...FieldAny)
}

// FieldAny is used to log any structed info
// {"key1": "value1", "key2": 2, "key3": [1,2,3]}
type FieldAny struct {
	key   string
	value interface{}
}

// Field returns a new FieldAny
func Field(key string, value interface{}) FieldAny {
	return FieldAny{
		key:   key,
		value: value,
	}
}

func (f FieldAny) zap() zap.Field {
	return zap.Any(f.key, f.value)
}

// Log implements logger interface
type Log struct {
	zap *zap.Logger
}

// New returns a new logger with default options
func New() *Log {
	loggerConf := zap.NewProductionConfig()
	loggerConf.DisableCaller = true
	loggerConf.DisableStacktrace = true
	loggerConf.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder

	logger, _ := loggerConf.Build()
	logger.WithOptions(
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
	)

	return &Log{
		zap: logger,
	}
}

// Info writes log message with INFO level.
func (l *Log) Info(message string, fields ...FieldAny) {
	l.zap.Info(message, convertFields(fields...)...)
}

// Error writes log message with ERROR level.
func (l *Log) Error(message string, fields ...FieldAny) {
	l.zap.Error(message, convertFields(fields...)...)
}

func convertFields(fields ...FieldAny) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))

	for i := range fields {
		zapFields = append(zapFields, fields[i].zap())
	}

	return zapFields
}
