package errors

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

// MaxStackTraceLimit allows to print only max 5 frames
const MaxStackTraceLimit = 5

// AppError struct holds the value of HTTP status code and custom error message.
// https://go.dev/blog/error-handling-and-go
type AppError struct {
	status       int         // `json:"status"` // HTTP status code
	message      string      // `json:"error,omitempty"`
	debug        error       // `json:"-"`
	conflictData interface{} // `json:"conflict_data,omitempty"` // Add any relevant trace info for debug
	errorDetails *Details    // `json:"error_details,omitempty"` // Custom internal error codes
}

// NewAppError returns the new apperror object
func NewAppError(status int, msg string) *AppError {
	return &AppError{
		status:  status,
		message: msg,
		debug:   ner(msg),
	}
}

// Error implements error interface
func (err *AppError) Error() string {
	return err.message
}

// GetStatus returns the status code of the apperror
func (err *AppError) GetStatus() int {
	return err.status
}

// UpdateStatus lets overrides the status code which is passed
func (err *AppError) UpdateStatus(statusCode int) {
	if statusCode == 0 || http.StatusText(statusCode) == "" {
		return
	}
	err.status = statusCode
}

// UpdateMsg lets overrides the error message will be send to the client
func (err *AppError) UpdateMsg(msg string) {
	err.message = msg
}

// MarshalJSON converts the AppError struct to JSON bytes
func (err *AppError) MarshalJSON() ([]byte, error) {
	if err == nil {
		return nil, ner("AppError: MarshalJSON on nil app error")
	}

	m := map[string]interface{}{
		"status": err.status,
		"error":  err.Error(), // err.Message
	}
	if err.conflictData != nil {
		m["conflict_data"] = err.conflictData
	}
	if err.errorDetails != nil {
		m["error_details"] = err.errorDetails
	}

	return json.Marshal(m)
}

// UnmarshalJSON convets a byte back to the AppError struct
func (err *AppError) UnmarshalJSON(b []byte) error {
	if err == nil {
		return ner("AppError: UnmarshalJSON on nil pointer")
	}

	var appErr struct {
		Status       int         `json:"status"`
		Message      string      `json:"error"`
		ConflictData interface{} `json:"conflict_data"`
		ErrorDetails *Details    `json:"error_details"`
	}
	if err := json.Unmarshal(b, &appErr); err != nil {
		return err
	}

	err.status = appErr.Status
	err.message = appErr.Message
	err.conflictData = appErr.ConflictData
	err.errorDetails = appErr.ErrorDetails
	return nil
}

// Using stackTracer, the stack trace of an error can be retrieved and controlled.
// https://pkg.go.dev/github.com/pkg/errors#hdr-Retrieving_the_stack_trace_of_an_error_or_wrapper
type stackTracer interface {
	StackTrace() errors.StackTrace
}

// Log logs the app error to stdout
// log stack trace
func (err *AppError) Log() {
	if err != nil {
		log.Println("[error-msg] ", err.status, err.Error())
		if err.debug == nil {
			return
		}
	}

	log.Println("[debug-error] " + err.debug.Error())
	stackTraceLimit := MaxStackTraceLimit
	if traceErr, ok := err.debug.(stackTracer); ok {
		if len(traceErr.StackTrace()) < stackTraceLimit {
			stackTraceLimit = len(traceErr.StackTrace())
		}

		log.Printf("[debug-error-trace]%+v\n", traceErr.StackTrace()[:stackTraceLimit])
	}
}

// GetDebug returns the debug error if presents else returns new error of the message
func (err *AppError) GetDebug() error {
	if err.debug != nil {
		return err.debug
	}

	return ner(err.message)
}

// AddDebug method is used to add a debug error which will be printed
// during the error execution if it is not nil. This is purely for developers
// debugging purposes
// This can be overridden N times in the entire lifecycle to wrap stacktraces
//
// expects cause of return type of wraped github.com/pkg/errors
func (err *AppError) AddDebug(cause error) *AppError {
	if err != nil {
		err.debug = cause
	}

	return err
}

// AddDebugf is a helper function which calls fmt.Errorf internally
// for the error which is added a Debug
func (err *AppError) AddDebugf(format string, a ...interface{}) *AppError {
	return err.AddDebug(fmt.Errorf(format, a...))
}

// DebugError returns the detailed warped error message if present, else returns default msg.
func (err *AppError) DebugError() string {
	if err.debug != nil {
		return err.debug.Error()
	}
	return err.message
}

// AddConflictData add extra conflict error information. This will be sent in meta
func (err *AppError) AddConflictData(data interface{}) *AppError {
	if err != nil {
		err.conflictData = data
	}

	return err
}

// GetConflictData returns the conflict data if present
func (err *AppError) GetConflictData() interface{} {
	return err.conflictData
}

// AddErrorDetail is used to add the error_details fields in the metadata
// object of the response. This field will be empty if ErrorDetail is not configured.
func (err *AppError) AddErrorDetail(d *Details) *AppError {
	err.errorDetails = d
	return err
}

// NotNil checks if the app errors is not nil or not
func (err *AppError) NotNil() bool {
	return err != nil
}

// IsBadRequest should return true if HTTP status of an error is 400.
func (err *AppError) IsBadRequest() bool {
	return err.status == http.StatusBadRequest
}

// IsForbidden should return true if HTTP status of an error is 403.
func (err *AppError) IsForbidden() bool {
	return err.status == http.StatusForbidden
}

// IsStatusNotFound should return true if HTTP status of an error is 404.
func (err *AppError) IsStatusNotFound() bool {
	return err.status == http.StatusNotFound
}

// IsInternalServerError should return true if HTTP status of an error is 500.
func (err *AppError) IsInternalServerError() bool {
	return err.status == http.StatusInternalServerError
}

// IsInternalError should return true if HTTP status of an error is >= 5XX.
func (err *AppError) IsInternalError() bool {
	return err.status >= http.StatusInternalServerError
}

const (
	requestCancelled      = "context canceled"
	statusClientCancelled = 460
)

// OverwriteStatusCode change the status code to 460, if the error containts
// context canceled.
func (err *AppError) OverwriteStatusCode() {
	if err == nil {
		return
	}
	if strings.Contains(err.Error(), requestCancelled) {
		err.status = statusClientCancelled
		return
	}

	// check in the debug error
	if err.debug != nil && strings.Contains(err.debug.Error(), requestCancelled) {
		err.status = statusClientCancelled
		return
	}
}
