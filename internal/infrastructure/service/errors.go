package service

// ErrorType represents type of service error
type ErrorType int

const (
	// BadRequest denotes that user input is incorrect
	BadRequest ErrorType = iota + 1
	// InternalError denotes infrastructure error
	InternalError
	// Unauthorized denotes that user is not authorized
	Unauthorized
	// NotFound speaks for itself
	NotFound
	// Forbidden means user has no access for a resource
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
