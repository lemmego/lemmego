package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/lemmego/api/app"
	"log"
	"net/http"
	"strings"
	"time"
)

type LogOptions struct {
	Headers   bool
	UserAgent bool
}

type responseRecorder struct {
	http.ResponseWriter
	status int
}

// WriteHeader wraps the ResponseWriter's WriteHeader method to capture the status code.
func (r *responseRecorder) WriteHeader(statusCode int) {
	r.status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

// Status method returns the status code.
func (r *responseRecorder) Status() int {
	return r.status
}

func LogRequest(opts ...*LogOptions) app.HTTPMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Capture the start time to measure latency
			start := time.Now()

			var logStr strings.Builder
			logStr.WriteString(fmt.Sprintf("[%s] <- %s %s", r.Method, r.URL, r.RemoteAddr))

			if len(opts) > 0 && opts[0].Headers {
				headers, err := json.Marshal(r.Header)
				if err != nil {
					panic(err)
				}
				logStr.WriteString(fmt.Sprintf("\n\nHeaders: %s\n\n", headers))
			}

			// Log request details
			log.Printf(logStr.String())

			// Create a responseRecorder to capture the status
			recorder := &responseRecorder{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(recorder, r)

			// Calculate the time taken for this request
			latency := time.Since(start)

			// Log after the request is processed including the status code
			log.Printf(
				"[%s] -> %s %s - %d | Completed in %v",
				r.Method,
				r.URL.Path,
				r.RemoteAddr,
				recorder.Status(),
				latency,
			)
		})
	}
}
