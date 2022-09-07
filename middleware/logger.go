package middleware

import (
	"log"
	"net/http"
	"time"
)

// RequestIDHeader is the name of the HTTP Header which contains the request id.
// Exported so that it can be changed by developers
var RequestIDHeader = "X-Request-Id"

// ResponseWriter hijacker, just hijack only status code

type responseWriter struct {
	http.ResponseWriter
	code int
}

// WriteHeader writes the status code in responsewriter
func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.code = statusCode

	rw.ResponseWriter.WriteHeader(statusCode)
}

func getReqID(r *http.Request) string {
	return r.Context().Value(RequestIDHeader).(string)
}

// Logger middlwware logs the request stats post the call.
// this hijacks the responsewrite status code.
// use logger middleware at the root of the router chain as possible
// TODO: logger might set the common log object in the ctx
func Logger(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{
			ResponseWriter: w,
		}

		// call the next handler
		next.ServeHTTP(rw, r)

		log.Printf(
			"[%s] %s [%d] %s?%s %v",
			getReqID(r), r.Method, rw.code, r.URL.Path, r.URL.RawQuery, time.Since(start),
		)
	}

	return http.HandlerFunc(fn)
}
