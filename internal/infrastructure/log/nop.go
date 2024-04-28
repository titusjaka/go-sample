package log

type NopLogger struct{}

func (NopLogger) Info(string, ...FieldAny) {}

func (NopLogger) Error(string, ...FieldAny) {}
