package httplogger

import (
	"io"
	"net/http"
	"time"
)

type private interface {
	private()
}

// ResponseLog is the information of the response.
type ResponseLog interface {
	// Header returns HTTP Header of the response.
	Header() http.Header

	// Status returns HTTP Status code.
	Status() int

	// The size of response body.
	//
	// Deprecated: Use ResponseSize instead.
	Size() int

	// Time returns the time when the request was received.
	//
	// Deprecated: Use RequestTime instead.
	Time() time.Time

	// RequestSize returns the size of request body.
	RequestSize() int64

	// ResponseSize returns the size of response body.
	ResponseSize() int64

	// RequestTime returns the time when the request was received.
	RequestTime() time.Time

	// WriteHeaderTime returns the time when the response header was written.
	WriteHeaderTime() time.Time

	private // avoid users to implement this interface
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
	req := *r // shallow copy
	body := &sizeReader{r: r.Body}
	req.Body = body
	lrw := &responseWriter{
		rw:          w,
		req:         &req,
		logger:      h.logger,
		requestTime: time.Now(),
	}
	h.handler.ServeHTTP(wrap(lrw), r)
	if !lrw.hijacked {
		io.Copy(io.Discard, body)
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
