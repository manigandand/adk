package respond

import (
	"encoding/json"
	"net/http"

	"github.com/manigandand/adk/errors"
	log "github.com/sirupsen/logrus"
)

// Response struct contains all the fields needed to respond
// to a particular request
// NOTE: we may add support for
// Metadata   interface{}
type Response struct {
	statusCode int
	data       interface{}
	headers    map[string]string
}

// sendResponse is a helper function which sends a JSON response with the passed data
func sendResponse(w http.ResponseWriter, statusCode int, data interface{}) *errors.AppError {
	if err := NewResponse(statusCode, data).Send(w); err != nil {
		log.Error("respond.send.error: ", err)
		// TODO: handle err, notify developers. Configurable
		// http.Error(w, err.Error(), http.StatusInternalServerError)

		return errors.InternalServer(err.Error()).AddDebug(err)
	}
	return nil
}

// NewResponse returns a new response object.
func NewResponse(statusCode int, data interface{}) *Response {
	return &Response{
		statusCode: statusCode,
		data:       data,
	}
}

// Send sends data encoded to JSON
func (res *Response) Send(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	if res.headers != nil {
		for key, value := range res.headers {
			w.Header().Set(key, value)
		}
	}

	w.WriteHeader(res.statusCode)
	if res.statusCode == http.StatusNoContent {
		return nil
	}

	return json.NewEncoder(w).Encode(res.data)
}

// 2xx JSON Response------------------------------------------------------------

// OK is a helper function used to send response data
// with StatusOK status code (200)
// NOTE: HTTP Methods GET, PUT should use status code (200)
func OK(w http.ResponseWriter, data interface{}) *errors.AppError {
	return sendResponse(w, http.StatusOK, data)
}

// Created is a helper function used to send response data
// with StatusCreated status code (201)
// NOTE: HTTP Method POST which creates any resource on the server
// should use status code (201)
func Created(w http.ResponseWriter, data interface{}) *errors.AppError {
	return sendResponse(w, http.StatusCreated, data)
}

// NoContent is a helper function used to send a NoContent Header (204)
// Note : the sent data and meta are ignored.
// NOTE: HTTP Methods DELETE should use status code (204)
func NoContent(w http.ResponseWriter, data interface{}) *errors.AppError {
	return sendResponse(w, http.StatusNoContent, nil)
}

// 4xx & 5XX JSON Response------------------------------------------------------

// Fail write the error response
// Common func to send all the error response
func Fail(w http.ResponseWriter, e *errors.AppError) *errors.AppError {
	e.Log()
	e.OverwriteStatusCode()

	return sendResponse(w, e.GetStatus(), e)
}

// 2xx CSV Response-------------------------------------------------------------

// CSV writes the text/csv response into the responsewriter
func CSV(w http.ResponseWriter, data string) *errors.AppError {
	w.Header().Add("Content-Type", "text/csv")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(data)); err != nil {
		return errors.InternalServer(err.Error()).AddDebug(err)
	}

	return nil
}
