package log

// contextLog represents log Entry. Slice of fields is used as context for each log message
type contextLog struct {
	logger Logger
	fields []FieldAny
}

func newContext(logger Logger) *contextLog {
	if c, ok := logger.(*contextLog); ok {
		return c
	}
	return &contextLog{logger: logger}
}

// Info writes log message with INFO level.
// It appends all contextual fields to input fields and delegates output to embedded logger
func (c *contextLog) Info(message string, fields ...FieldAny) {
	c.logger.Info(message, append(c.fields, fields...)...)
}

// Error writes log message with ERROR level.
// It appends all contextual fields to input fields and delegates output to embedded logger
func (c *contextLog) Error(message string, fields ...FieldAny) {
	c.logger.Error(message, append(c.fields, fields...)...)
}

// With returns logger with input fields as context.
// If input Logger is NOT *contextLog, a new *contextLog is created and fields are passed to *contextLog.
// If input Logger is *contextLog, input fields are appended to existing logger fields
func With(logger Logger, fields ...FieldAny) Logger {
	if len(fields) == 0 {
		return logger
	}
	l := newContext(logger)
	return &contextLog{
		logger: l.logger,
		fields: append(l.fields, fields...),
	}
}
