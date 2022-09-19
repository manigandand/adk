package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

// ContextKey avoids collision in Context.
// Read more on https://blog.golang.org/context#TOC_3.2.
type ContextKey string

func (c ContextKey) String() string {
	return string(c)
}

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
	reqID, ok := r.Context().Value(RequestIDHeader).(string)
	if !ok {
		return ""
	}
	return reqID
}

// Logger middlwware logs the request stats post the call.
// this hijacks the responsewrite status code.
// use logger middleware at the root of the router chain as possible
// TODO: logger might set the common log object in the ctx
func Logger(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		reqID := r.Header.Get(RequestIDHeader)
		if reqID == "" {
			reqID = getReqID(r)
			if reqID == "" {
				reqID = fmt.Sprintf("%d", time.Now().UnixNano())
				ctx = context.WithValue(ctx, ContextKey(RequestIDHeader), reqID)
			}
		}
		start := time.Now()
		rw := &responseWriter{
			ResponseWriter: w,
		}

		// call the next handler
		next.ServeHTTP(rw, r.WithContext(ctx))

		log.Printf(
			"[%s] %s [%d] %s?%s %v",
			reqID, r.Method, rw.code, r.URL.Path, r.URL.RawQuery, time.Since(start),
		)
	}

	return http.HandlerFunc(fn)
}
