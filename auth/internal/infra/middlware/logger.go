package middleware

import (
	"log"
	"net/http"
	"time"
)

func Logger() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			lw := &responseLogger{w: w, code: http.StatusOK}

			defer func() {
				log.Printf(
					"method=%s path=%s status=%d duration=%s",
					r.Method, r.URL.Path, lw.code, time.Since(start),
				)
			}()

			next.ServeHTTP(lw, r)
		})
	}
}

type responseLogger struct {
	w    http.ResponseWriter
	code int
}

func (l *responseLogger) Header() http.Header {
	return l.w.Header()
}

func (l *responseLogger) Write(b []byte) (int, error) {
	return l.w.Write(b)
}

func (l *responseLogger) WriteHeader(code int) {
	l.code = code
	l.w.WriteHeader(code)
}
