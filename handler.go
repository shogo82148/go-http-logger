package httplogger

import (
	"net/http"
	"time"
)

type ResponseLog interface {
	Header() http.Header // HTTP Header
	Status() int         // HTTP Status code
	Size() int           // The size of response body
	Time() time.Time     // Time the request was received
}

type Logger interface {
	WriteHTTPLog(l ResponseLog, r *http.Request)
}

func LoggingHandler(logger Logger, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := makeLogger(w)
		handler.ServeHTTP(l, r)
		logger.WriteHTTPLog(l, r)
	})
}

type LoggerFunc func(l ResponseLog, r *http.Request)

func (f LoggerFunc) WriteHTTPLog(l ResponseLog, r *http.Request) {
	f(l, r)
}
