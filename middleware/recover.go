package middleware

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/manigandand/adk/errors"
	"github.com/manigandand/adk/respond"
	log "github.com/sirupsen/logrus"
)

// Recoverer is a middleware that recovers from panics, logs the panic (and a
// backtrace), and returns a HTTP 500 (Internal Server Error) status if
// possible. Recoverer prints a request ID if one is provided.
func Recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				const size = 4096

				buf := make([]byte, size)
				buf = buf[:runtime.Stack(buf, false)]
				buffer := string(buf)

				excMessage := fmt.Errorf("%+v", rvr).Error()

				log.Error(fmt.Errorf("[panic-recover] %v \n%v", excMessage, buffer), r.URL.String())

				// TODO: Let devloper know the panic
				// Log into sentry/..
				respond.Fail(w, errors.InternalServerStd())
				return
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
