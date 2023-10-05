package errors

import (
	"net/http"

	"github.com/pkg/errors"
)

var ner = errors.New

// New github.com/pkg/errors wraper
// New returns an error with the supplied message.
// New also records the stack trace at the point it was called.
func New(message string) error {
	return ner(message)
}

// Errorf github.com/pkg/errors wraper
// Errorf formats according to a format specifier and returns the string
// as a value that satisfies error.
// Errorf also records the stack trace at the point it was called.
func Errorf(format string, args ...interface{}) error {
	return errors.Errorf(format, args...)
}

// Wrapf github.com/pkg/errors.Wrapf
// Wrapf returns an error annotating err with a stack trace
// at the point Wrapf is called, and the format specifier.
// If err is nil, Wrapf returns nil.
func Wrapf(err error, format string, args ...interface{}) error {
	return errors.Wrapf(err, format, args...)
}

// Wrap github.com/pkg/errors.Wrap
// Wrap returns an error annotating err with a stack trace
// at the point Wrap is called, and the supplied message.
// If err is nil, Wrap returns nil.
func Wrap(err error, message string) error {
	return errors.Wrap(err, message)
}

// KeyRequired returns new error with custom error message
func KeyRequired(key string) *AppError {
	return BadRequest(key + " is required")
}

// InvalidKey returns new error with custom error message
func InvalidKey(val, key string) *AppError {
	return BadRequest(val + " is invalid " + key)
}

// 4xx -------------------------------------------------------------------------

// BadRequest will return `http.StatusBadRequest` with custom message.
func BadRequest(message string) *AppError { // 400
	return NewAppError(http.StatusBadRequest, message)
}

// Unauthorized will return `http.StatusUnauthorized` with custom message.
func Unauthorized(message string) *AppError { // 401
	return NewAppError(http.StatusUnauthorized, message)
}

// PaymentRequired will return `http.StatusPaymentRequired` with custom message.
func PaymentRequired(message string) *AppError { // 402
	return NewAppError(http.StatusPaymentRequired, message)
}

// Forbidden will return `http.StatusForbidden` with custom message.
func Forbidden(message string) *AppError { // 403
	return NewAppError(http.StatusForbidden, message)
}

// NotFound will return `http.StatusNotFound` with custom message.
func NotFound(message string) *AppError { // 404
	return NewAppError(http.StatusNotFound, message)
}

// Conflict will return `http.StatusConflict` with custom message.
func Conflict(message string) *AppError { // 409
	return NewAppError(http.StatusConflict, message)
}

// Gone will return `http.StatusGone` with custom message.
func Gone(message string) *AppError { // 410
	return NewAppError(http.StatusGone, message)
}

// UnprocessableEntity will return `http.StatusUnprocessableEntity` with
// custom message.
func UnprocessableEntity(message string) *AppError { // 422
	return NewAppError(http.StatusUnprocessableEntity, message)
}

// TooManyRequests will return `http.StatusTooManyRequests` with
// custom message.
func TooManyRequests(message string) *AppError { // 422
	return NewAppError(http.StatusTooManyRequests, message)
}

// TooEarly will return `http.StatusTooEarly` with
// custom message.
func TooEarly(message string) *AppError { // 425
	return NewAppError(http.StatusTooEarly, message)
}

// 5xx -------------------------------------------------------------------------

// InternalServer will return `http.StatusInternalServerError` with custom message.
func InternalServer(message string) *AppError { // 500
	return NewAppError(http.StatusInternalServerError, message)
}

// InternalServerStd will return `http.StatusInternalServerError` with static message.
func InternalServerStd() *AppError { // 500
	return NewAppError(http.StatusInternalServerError, "something went wrong")
}
