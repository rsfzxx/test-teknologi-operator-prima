package apperror

import "net/http"

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type Error struct {
	Status  int
	Message string
	Details []FieldError
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *Error) Unwrap() error { return e.Err }

func New(status int, message string, details ...FieldError) *Error {
	return &Error{Status: status, Message: message, Details: details}
}

func BadRequest(message string, details ...FieldError) *Error {
	return New(http.StatusBadRequest, message, details...)
}

func NotFound(message string) *Error {
	return New(http.StatusNotFound, message)
}

func MethodNotAllowed() *Error {
	return New(http.StatusMethodNotAllowed, "Method not allowed")
}

func Validation(details ...FieldError) *Error {
	return New(http.StatusUnprocessableEntity, "Validation failed", details...)
}

func Internal(err error) *Error {
	return &Error{
		Status:  http.StatusInternalServerError,
		Message: "Internal server error",
		Err:     err,
	}
}

func Conflict(message string, details ...FieldError) *Error {
	return New(http.StatusConflict, message, details...)
}