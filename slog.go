//go:build go1.21
// +build go1.21

package httplogger

import (
	"log/slog"
	"net/http"
)

type slogLogger struct {
	level  slog.Level
	msg    string
	logger *slog.Logger
}

// NewSlogLogger returns a Logger that logs to the given slog.Logger.
// Logged attributes are:
//
//   - received_time: Time the request was received.
//   - host: Remote host.
//   - forwardedfor: X-Forwarded-For header.
//   - user: Remote user.
//   - req: First line of request.
//   - method: Request method.
//   - uri: Request URI.
//   - protocol: Requested Protocol (usually "HTTP/1.0" or "HTTP/1.1").
//   - status: HTTP status code of the response.
//   - sent_bytes: Size of response body in bytes, excluding HTTP headers.
//   - received_bytes: Size of request body in bytes, excluding HTTP headers.
//   - referer: Referer header.
//   - ua: User agent header.
//   - request_time: Time taken to serve the request, in seconds.
func NewSlogLogger(level slog.Level, msg string, logger *slog.Logger) Logger {
	return &slogLogger{
		level:  level,
		msg:    msg,
		logger: logger,
	}
}

// WriteHTTPLog implements the Logger interface.
func (log *slogLogger) WriteHTTPLog(attrs Attrs, r *http.Request) {
	reqtime := attrs.WriteHeaderTime().Sub(attrs.RequestTime())
	user, _, _ := r.BasicAuth()

	// the names of attrs are based on https://github.com/tkuchiki/alp/blob/main/README.md#json
	// some of them come from https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-access-logs.html
	// and http://ltsv.org/.
	log.logger.Log(
		r.Context(),
		log.level,
		log.msg,

		// Time the request was received.
		//
		// "time" is most commonly used, but it is reserved by slog.
		// So, we use "received_time" instead.
		slog.Time("received_time", attrs.RequestTime()),

		// Remote host.
		slog.String("host", r.RemoteAddr),

		// X-Forwarded-For header.
		slog.String("forwardedfor", r.Header.Get("X-Forwarded-For")),

		// Remote user.
		slog.String("user", user),

		// First line of request.
		slog.String("req", r.Method+" "+r.RequestURI+" "+r.Proto),

		// Request method.
		slog.String("method", r.Method),

		// Request URI.
		slog.String("uri", r.RequestURI),

		// Requested Protocol (usually "HTTP/1.0" or "HTTP/1.1").
		slog.String("protocol", r.Proto),

		// Status code.
		slog.Int("status", attrs.Status()),

		// Size of response body in bytes, excluding HTTP headers.
		//
		// The key name "sent_bytes" comes from ALB log.
		// Note that ALB log includes HTTP headers however this value does not include them.
		slog.Int64("sent_bytes", attrs.ResponseSize()),

		// Size of request body in bytes, excluding HTTP headers.
		//
		// The key name "received_bytes" comes from ALB log.
		// Note that ALB log includes HTTP headers however this value does not include them.
		slog.Int64("received_bytes", attrs.RequestSize()),

		// Referer header.
		slog.String("referer", r.Referer()),

		// User-Agent header.
		slog.String("ua", r.UserAgent()),

		// Host header.
		slog.String("vhost", r.Host),

		// The time taken to serve the request in seconds.
		slog.Float64("request_time", reqtime.Seconds()),
	)
}
