//go:build go1.21
// +build go1.21

package httplogger

import (
	"log/slog"
)

type slogLogger struct {
	level  slog.Level
	msg    string
	logger slog.Logger
}

func NewSlogLogger(level slog.Level, msg string, logger slog.Logger) Logger {
	return &slogLogger{
		level:  level,
		msg:    msg,
		logger: logger,
	}
}

// WriteHTTPLog implements the Logger interface.
func (log *slogLogger) WriteHTTPLog(attrs Attrs, r *http.Request) {
	// the names of attrs come from http://ltsv.org/

	now := time.Now()
	reqtime := now.Sub(attrs.RequestTime())
	apptime := attrs.WriteHeaderTime().Sub(attrs.RequestTime())
	user, _, _ := r.BasicAuth()
	log.logger.Log(
		r.Context(),
		log.level,
		log.msg,

		// Time the request was received
		slog.Time("time", attrs.RequestSize()),

		// Remote host
		slog.String("host", r.RemoteAddr),

		// X-Forwarded-For header
		slog.String("forwardedfor", r.Header.Get("X-Forwarded-For")),

		// Remote user
		slog.String("user", user),

		// First line of request
		slog.String("req", r.Method+" "+r.RequestURI+" "+r.Proto),

		// Request method
		slog.String("method", r.Method),

		// Request URI
		slog.String("uri", r.RequestURI),

		// Requested Protocol (usually "HTTP/1.0" or "HTTP/1.1")
		slog.String("protocol", r.Proto),

		// Status code
		slog.Int("status", attrs.Status()),

		// Size of response in bytes, excluding HTTP headers.
		slog.Int64("size", attrs.ResponseSize()),

		// Size of response in bytes, excluding HTTP headers.
		slog.Int64("reqsize", attrs.RequestSize()),

		// Referer header
		slog.String("referer", r.Referer()),

		// User-Agent header
		slog.String("ua", r.UserAgent()),

		// Host header
		slog.String("vhost", r.Host),

		// The time taken to serve the request
		slog.Duration("reqtime", reqtime),

		// X-Cache header
		slog.String("cache", attrs.Header().Get("X-Cache")),

		// Execution time for processing some request, e.g. X-Runtime header for application server or processing time of SQL for DB server.
		slog.String("runtime", attrs.Header().Get("X-Runtime")),

		// Response time from the upstream server
		slog.Duration("apptime", apptime),
	)
}
