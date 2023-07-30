package httplogger

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
)

type dummyFlusher struct {
	http.ResponseWriter
	called bool
}

func (rw *dummyFlusher) Flush() {
	rw.called = true
}

func TestWrap_Flusher(t *testing.T) {
	rw := &dummyFlusher{}
	got := wrap(&responseWriter{rw: rw})
	if flusher, ok := got.(http.Flusher); !ok {
		t.Error("want to implement http.Flusher, but it doesn't")
	} else {
		flusher.Flush()
	}
	if !rw.called {
		t.Error("Flush() is not called")
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
	if _, ok := got.(io.StringWriter); ok {
		t.Error("want not to implement io.io.StringWriter, but it does")
	}
	if _, ok := got.(http.Pusher); ok {
		t.Error("want not to implement http.Pusher, but it does")
	}
}

type dummyCloseNotifier struct {
	http.ResponseWriter
	called bool
}

func (rw *dummyCloseNotifier) CloseNotify() <-chan bool {
	rw.called = true
	return nil
}

func TestWrap_CloseNotify(t *testing.T) {
	rw := &dummyCloseNotifier{}
	got := wrap(&responseWriter{rw: rw})
	if _, ok := got.(http.Flusher); ok {
		t.Error("want not to implement http.Flusher, but it does")
	}
	if notifier, ok := got.(http.CloseNotifier); !ok {
		t.Error("want to implement http.CloseNotifier, but it doesn't")
	} else {
		notifier.CloseNotify()
	}
	if !rw.called {
		t.Error("CloseNotify() is not called")
	}
	if _, ok := got.(http.Hijacker); ok {
		t.Error("want not to implement http.Hijacker, but it does")
	}
	if _, ok := got.(io.ReaderFrom); ok {
		t.Error("want not to implement http.ReaderFrom, but it does")
	}
	if _, ok := got.(io.StringWriter); ok {
		t.Error("want not to implement io.io.StringWriter, but it does")
	}
	if _, ok := got.(http.Pusher); ok {
		t.Error("want not to implement http.Pusher, but it does")
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
	if _, ok := got.(io.StringWriter); ok {
		t.Error("want not to implement io.io.StringWriter, but it does")
	}
	if _, ok := got.(http.Pusher); ok {
		t.Error("want not to implement http.Pusher, but it does")
	}
}

type dummyReaderFrom struct {
	http.ResponseWriter
	buf bytes.Buffer
}

func (rw *dummyReaderFrom) ReadFrom(r io.Reader) (n int64, err error) {
	return rw.buf.ReadFrom(r)
}

func TestWrap_ReaderFrom(t *testing.T) {
	rw := &dummyReaderFrom{}
	got := wrap(&responseWriter{rw: rw})
	if _, ok := got.(http.Flusher); ok {
		t.Error("want not to implement http.Flusher, but it does")
	}
	if _, ok := got.(http.CloseNotifier); ok {
		t.Error("want not to implement http.CloseNotifier, but it does")
	}
	if _, ok := got.(http.Hijacker); ok {
		t.Error("want not to implement http.Hijacker, but it does")
	}
	if reader, ok := got.(io.ReaderFrom); !ok {
		t.Error("want to implement http.ReaderFrom, but it doesn't")
	} else {
		reader.ReadFrom(strings.NewReader("hello"))
	}
	if rw.buf.String() != "hello" {
		t.Errorf("got %q, want %q", rw.buf.String(), "hello")
	}
	if _, ok := got.(io.StringWriter); ok {
		t.Error("want not to implement io.io.StringWriter, but it does")
	}
	if _, ok := got.(http.Pusher); ok {
		t.Error("want not to implement http.Pusher, but it does")
	}
}

type dummyStringWriter struct {
	http.ResponseWriter
	buf bytes.Buffer
}

func (rw *dummyStringWriter) WriteString(s string) (n int, err error) {
	return rw.buf.WriteString(s)
}

func TestWrap_StringWriter(t *testing.T) {
	rw := &dummyStringWriter{}
	got := wrap(&responseWriter{rw: rw})
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
	if sw, ok := got.(io.StringWriter); !ok {
		t.Error("want to implement io.io.StringWriter, but it doesn't")
	} else {
		sw.WriteString("hello")
	}
	if rw.buf.String() != "hello" {
		t.Errorf("want %q, but %q", "hello", rw.buf.String())
	}
	if _, ok := got.(http.Pusher); ok {
		t.Error("want not to implement http.Pusher, but it does")
	}
}

type dummyPusher struct {
	http.ResponseWriter
	called bool
}

func (rw *dummyPusher) Push(target string, opts *http.PushOptions) error {
	rw.called = true
	return nil
}

func TestWrap_Pusher(t *testing.T) {
	rw := &dummyPusher{}
	got := wrap(&responseWriter{rw: rw})
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
	if _, ok := got.(io.StringWriter); ok {
		t.Error("want not to implement io.io.StringWriter, but it does")
	}
	if pusher, ok := got.(http.Pusher); !ok {
		t.Error("want to implement http.Pusher, but it doesn't")
	} else {
		pusher.Push("", nil)
	}
	if !rw.called {
		t.Error("Push() is not called")
	}
}
