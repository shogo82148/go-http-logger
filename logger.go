package httplogger

//go:generate go run codegen.go

import (
	"bufio"
	"net"
	"net/http"
	"time"
)

// responseWriter is wrapper of http.ResponseWriter that keeps track of its HTTP
// status code and body size
type responseWriter struct {
	rw     http.ResponseWriter
	status int
	size   int
	t      time.Time
}

func (rw *responseWriter) Header() http.Header {
	return rw.rw.Header()
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.status == 0 {
		// The status will be StatusOK if WriteHeader has not been called yet
		rw.status = http.StatusOK
	}
	size, err := rw.rw.Write(b)
	rw.size += size
	return size, err
}

func (rw *responseWriter) WriteHeader(s int) {
	rw.rw.WriteHeader(s)
	rw.status = s
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
	if err == nil && rw.status == 0 {
		// The status will be StatusSwitchingProtocols if there was no error and
		// WriteHeader has not been called yet
		rw.status = http.StatusSwitchingProtocols
	}
	return conn, buf, err
}

func (rw *responseWriter) CloseNotify() <-chan bool {
	n := rw.rw.(http.CloseNotifier)
	return n.CloseNotify()
}

func (rw *responseWriter) Push(target string, opts *http.PushOptions) error {
	if p, ok := rw.rw.(http.Pusher); ok {
		return p.Push(target, opts)
	}
	return http.ErrNotSupported
}
