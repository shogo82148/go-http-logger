package httplogger

import (
	"net/http"
	"time"
)

// ResponseLog is the information of the response.
type ResponseLog interface {
	Header() http.Header // HTTP Header
	Status() int         // HTTP Status code
	Size() int           // The size of response body
	Time() time.Time     // Time the request was received
}

// Logger is the interface for your custom logger.
type Logger interface {
	WriteHTTPLog(l ResponseLog, r *http.Request)
}

type loggingHandler struct {
	logger  Logger
	handler http.Handler
}

func (h *loggingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lrw := &responseWriter{
		rw:     w,
		req:    r,
		logger: h.logger,
	}
	h.handler.ServeHTTP(wrap(lrw), r)
	if !lrw.hijacked {
		h.logger.WriteHTTPLog(lrw, r)
	}
}

// LoggingHandler wraps the HTTP handler with the logger.
func LoggingHandler(logger Logger, handler http.Handler) http.Handler {
	return &loggingHandler{
		logger:  logger,
		handler: handler,
	}
}

// The LoggerFunc type is an adapter to allow the use of ordinary functions as Logger.
type LoggerFunc func(l ResponseLog, r *http.Request)

// WriteHTTPLog implements the Logger interface.
func (f LoggerFunc) WriteHTTPLog(l ResponseLog, r *http.Request) {
	f(l, r)
}
