package log

import (
	"fmt"
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Level is just an alias for zapcore.Level and is used to set level for StdLogger
type Level = zapcore.Level

const (
	// Error level
	Error = zapcore.ErrorLevel
)

// NewStdLogger creates a new standard logger from Logger interface
// logger — must be either *Log or *contextLog
func NewStdLogger(logger Logger, level Level) (*log.Logger, error) {
	logStruct, err := newLogFromLogger(logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create std logger: %w", err)
	}
	return zap.NewStdLogAt(logStruct.zap, level)
}

func newLogFromLogger(logger Logger) (*Log, error) {
	switch logger := logger.(type) {
	case *Log:
		return logger, nil
	case *contextLog:
		return newLogFromLogger(logger.logger)
	default:
		return nil, fmt.Errorf("failed to retrieve logger, unsupported type %T", logger)
	}
}
