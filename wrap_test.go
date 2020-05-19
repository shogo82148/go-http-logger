package httplogger

import (
	"bufio"
	"io"
	"net"
	"net/http"
	"testing"
)

type dummyFlusher struct {
	http.ResponseWriter
}

func (dummyFlusher) Flush() {
	panic("unreachable")
}

func TestWrap_Flusher(t *testing.T) {
	got := wrap(&responseWriter{rw: dummyFlusher{}})
	if _, ok := got.(http.Flusher); !ok {
		t.Error("want to implement http.Flusher, but it doesn't")
	}
	if _, ok := got.(http.CloseNotifier); ok {
		t.Error("want not to implement http.CloseNotifier, but it does")
	}
	if _, ok := got.(http.Hijacker); ok {
		t.Error("want not to implement http.Hijacker, but it does")
	}
	if _, ok := got.(io.ReaderFrom); ok {
		t.Error("want not to implement http.ReaderFrom, but it does")
	}
	if _, ok := got.(stringWriter); ok {
		t.Error("want not to implement io.StringWriter, but it does")
	}
}

type dummyCloseNotifier struct {
	http.ResponseWriter
}

func (dummyCloseNotifier) CloseNotify() <-chan bool {
	panic("unreachable")
}

func TestWrap_CloseNotify(t *testing.T) {
	got := wrap(&responseWriter{rw: dummyCloseNotifier{}})
	if _, ok := got.(http.Flusher); ok {
		t.Error("want not to implement http.Flusher, but it does")
	}
	if _, ok := got.(http.CloseNotifier); !ok {
		t.Error("want to implement http.CloseNotifier, but it doesn't")
	}
	if _, ok := got.(http.Hijacker); ok {
		t.Error("want not to implement http.Hijacker, but it does")
	}
	if _, ok := got.(io.ReaderFrom); ok {
		t.Error("want not to implement http.ReaderFrom, but it does")
	}
	if _, ok := got.(stringWriter); ok {
		t.Error("want not to implement io.StringWriter, but it does")
	}
}

type dummyHijacker struct {
	http.ResponseWriter
}

func (dummyHijacker) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	panic("unreachable")
}

func TestWrap_Hijacker(t *testing.T) {
	got := wrap(&responseWriter{rw: dummyHijacker{}})
	if _, ok := got.(http.Flusher); ok {
		t.Error("want not to implement http.Flusher, but it does")
	}
	if _, ok := got.(http.CloseNotifier); ok {
		t.Error("want not to implement http.CloseNotifier, but it does")
	}
	if _, ok := got.(http.Hijacker); !ok {
		t.Error("want to implement http.Hijacker, but it doesn't")
	}
	if _, ok := got.(io.ReaderFrom); ok {
		t.Error("want not to implement http.ReaderFrom, but it does")
	}
	if _, ok := got.(stringWriter); ok {
		t.Error("want not to implement io.StringWriter, but it does")
	}
}

type dummyReaderFrom struct {
	http.ResponseWriter
}

func (dummyReaderFrom) ReadFrom(r io.Reader) (n int64, err error) {
	panic("unreachable")
}

func TestWrap_ReaderFrom(t *testing.T) {
	got := wrap(&responseWriter{rw: dummyReaderFrom{}})
	if _, ok := got.(http.Flusher); ok {
		t.Error("want not to implement http.Flusher, but it does")
	}
	if _, ok := got.(http.CloseNotifier); ok {
		t.Error("want not to implement http.CloseNotifier, but it does")
	}
	if _, ok := got.(http.Hijacker); ok {
		t.Error("want not to implement http.Hijacker, but it does")
	}
	if _, ok := got.(io.ReaderFrom); ok {
		t.Error("want to implement http.ReaderFrom, but it doesn't")
	}
	if _, ok := got.(stringWriter); ok {
		t.Error("want not to implement io.StringWriter, but it does")
	}
}

type dummyStringWriter struct {
	http.ResponseWriter
}

func (dummyStringWriter) WriteString(s string) (n int, err error) {
	panic("unreachable")
}

func TestWrap_StringWriter(t *testing.T) {
	got := wrap(&responseWriter{rw: dummyStringWriter{}})
	if _, ok := got.(http.Flusher); ok {
		t.Error("want not to implement http.Flusher, but it does")
	}
	if _, ok := got.(http.CloseNotifier); ok {
		t.Error("want not to implement http.CloseNotifier, but it does")
	}
	if _, ok := got.(http.Hijacker); ok {
		t.Error("want not to implement http.Hijacker, but it does")
	}
	if _, ok := got.(io.ReaderFrom); ok {
		t.Error("want not to implement http.ReaderFrom, but it does")
	}
	if _, ok := got.(stringWriter); ok {
		t.Error("want to implement io.StringWriter, but it doesn't")
	}
}
