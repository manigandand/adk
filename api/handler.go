package api

import (
	"log"
	"net/http"

	"github.com/manigandand/adk/errors"
	"github.com/manigandand/adk/respond"
)

// API Handler's ---------------------------------------------------------------

// Handler custom api handler help us to handle all the errors in one place
type Handler func(w http.ResponseWriter, r *http.Request) *errors.AppError

// ServeHTTP implements http handler interface
func (fn Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// NOTE: this expects the hijacked responsewriter to catch the statuscode
	// that will be logged in logger middleware.
	// use logger middleware at the root of the router chain as possible

	if err := fn(w, r); err != nil {
		// TODO: handle 5XX, notify developers. Configurable
		if fErr := respond.Fail(w, err); fErr.NotNil() {
			log.Printf("[panic] failed to write response. [%s] %s [%d] %s?%s",
				getReqID(r), r.Method, err.GetStatus(), r.URL.Path, r.URL.RawQuery,
			)
		}
	}
}

// RequestIDHeader is the name of the HTTP Header which contains the request id.
// Exported so that it can be changed by developers
var RequestIDHeader = "X-Request-Id"

func getReqID(r *http.Request) string {
	return r.Context().Value(RequestIDHeader).(string)
}
