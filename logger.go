package httplogger

//go:generate go run codegen.go

import (
	"bufio"
	"io"
	"log"
	"net"
	"net/http"
	"path"
	"runtime"
	"strings"
	"time"
)

// backport of io.StringWriter from Go 1.11
type stringWriter interface {
	WriteString(s string) (n int, err error)
}

type rwUnwrapper interface {
	Unwrap() http.ResponseWriter
}

// responseWriter is wrapper of http.ResponseWriter that keeps track of its HTTP
// status code and body size
type responseWriter struct {
	rw          http.ResponseWriter
	req         *http.Request
	logger      Logger
	wroteHeader bool
	status      int
	size        int
	t           time.Time
	hijacked    bool
}

func (rw *responseWriter) Header() http.Header {
	return rw.rw.Header()
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.wroteHeader {
		// The status will be StatusOK if WriteHeader has not been called yet
		rw.rw.WriteHeader(http.StatusOK)
		rw.status = http.StatusOK
		rw.wroteHeader = true
	}
	size, err := rw.rw.Write(b)
	rw.size += size
	return size, err
}

func (rw *responseWriter) WriteHeader(s int) {
	if rw.wroteHeader {
		caller := relevantCaller()
		log.Printf("httplogger: superfluous response.WriteHeader call from %s (%s:%d)", caller.Function, path.Base(caller.File), caller.Line)
		return
	}
	rw.rw.WriteHeader(s)
	rw.status = s
	rw.wroteHeader = true
}

// relevantCaller searches the call stack for the first function outside of net/http.
// The purpose of this function is to provide more helpful error messages.
func relevantCaller() runtime.Frame {
	pc := make([]uintptr, 16)
	n := runtime.Callers(1, pc)
	frames := runtime.CallersFrames(pc[:n])
	var frame runtime.Frame
	for {
		frame, more := frames.Next()
		if !strings.HasPrefix(frame.Function, "github.com/shogo82148/go-http-logger.") {
			return frame
		}
		if !more {
			break
		}
	}
	return frame
}

func (rw *responseWriter) Status() int {
	if rw.status == 0 {
		// The status will be StatusOK if WriteHeader has not been called yet
		rw.status = http.StatusOK
	}
	return rw.status
}

func (rw *responseWriter) Size() int {
	return rw.size
}

func (rw *responseWriter) Flush() {
	f, ok := rw.rw.(http.Flusher)
	if ok {
		f.Flush()
	}
}

func (rw *responseWriter) Time() time.Time {
	return rw.t
}

func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h := rw.rw.(http.Hijacker)
	conn, buf, err := h.Hijack()
	if err == nil {
		// The status will be StatusSwitchingProtocols if there was no error and
		// WriteHeader has not been called yet
		rw.status = http.StatusSwitchingProtocols
		rw.hijacked = true
		rw.wroteHeader = true
		rw.logger.WriteHTTPLog(rw, rw.req)
	}
	return conn, buf, err
}

func (rw *responseWriter) CloseNotify() <-chan bool {
	n := rw.rw.(http.CloseNotifier)
	return n.CloseNotify()
}

func (rw *responseWriter) WriteString(str string) (int, error) {
	if s, ok := rw.rw.(stringWriter); ok {
		return s.WriteString(str)
	}
	return rw.rw.Write([]byte(str))
}

func (rw *responseWriter) ReadFrom(src io.Reader) (n int64, err error) {
	if r, ok := rw.rw.(io.ReaderFrom); ok {
		return r.ReadFrom(src)
	}
	return io.Copy(rw.rw, src)
}

func (rw *responseWriter) Push(target string, opts *http.PushOptions) error {
	if p, ok := rw.rw.(http.Pusher); ok {
		return p.Push(target, opts)
	}
	return http.ErrNotSupported
}

// Unwrap returns the original http.ResponseWriter underlying this.
// It is used by [net/http.ResponseController].
func (rw *responseWriter) Unwrap() http.ResponseWriter {
	return rw.rw
}
