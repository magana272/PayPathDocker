package middleware

import (
	"net/http"
	"time"

	"paypath/pkg/logger"
)

type statusWriter struct {
	http.ResponseWriter
	code int
}

func (w *statusWriter) WriteHeader(code int) {
	w.code = code
	w.ResponseWriter.WriteHeader(code)
}

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &statusWriter{ResponseWriter: w, code: 200}
		next.ServeHTTP(sw, r)
		logger.Log.Info().
			Str("method", r.Method).
			Str("path", r.URL.RequestURI()).
			Int("status", sw.code).
			Dur("latency", time.Since(start)).
			Msg("request")
	})
}
