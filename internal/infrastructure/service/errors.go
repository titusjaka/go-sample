package service

// ErrorType represents type of service error
type ErrorType int

const (
	BadRequest ErrorType = iota + 1
	InternalError
	Unauthorized
	NotFound
	Forbidden
)

// Error error represents any business or infrastructure error
type Error struct {
	Type ErrorType
	Base error
}

// Unwrap implements errors.Wrapper interface
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Base
}

// Error implements errors.Error interface
func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	return e.Base.Error()
}
