// Code generated by shogo82148/go-http-logger/codegen.go; DO NOT EDIT

// +build !go1.8

package httplogger

import (
	"io"
	"net/http"
)

func wrap(rw *responseWriter) http.ResponseWriter {
	var n uint
	if _, ok := rw.rw.(http.Flusher); ok {
		n |= 0x1
	}
	if _, ok := rw.rw.(http.CloseNotifier); ok {
		n |= 0x2
	}
	if _, ok := rw.rw.(http.Hijacker); ok {
		n |= 0x4
	}
	if _, ok := rw.rw.(io.ReaderFrom); ok {
		n |= 0x8
	}
	if _, ok := rw.rw.(stringWriter); ok {
		n |= 0x10
	}
	switch n {
	case 0x0:
		return struct {
			http.ResponseWriter
		}{rw}
	case 0x1:
		return struct {
			http.ResponseWriter
			stringWriter
		}{rw, rw}
	case 0x2:
		return struct {
			http.ResponseWriter
			io.ReaderFrom
		}{rw, rw}
	case 0x3:
		return struct {
			http.ResponseWriter
			io.ReaderFrom
			stringWriter
		}{rw, rw, rw}
	case 0x4:
		return struct {
			http.ResponseWriter
			http.Hijacker
		}{rw, rw}
	case 0x5:
		return struct {
			http.ResponseWriter
			http.Hijacker
			stringWriter
		}{rw, rw, rw}
	case 0x6:
		return struct {
			http.ResponseWriter
			http.Hijacker
			io.ReaderFrom
		}{rw, rw, rw}
	case 0x7:
		return struct {
			http.ResponseWriter
			http.Hijacker
			io.ReaderFrom
			stringWriter
		}{rw, rw, rw, rw}
	case 0x8:
		return struct {
			http.ResponseWriter
			http.CloseNotifier
		}{rw, rw}
	case 0x9:
		return struct {
			http.ResponseWriter
			http.CloseNotifier
			stringWriter
		}{rw, rw, rw}
	case 0xa:
		return struct {
			http.ResponseWriter
			http.CloseNotifier
			io.ReaderFrom
		}{rw, rw, rw}
	case 0xb:
		return struct {
			http.ResponseWriter
			http.CloseNotifier
			io.ReaderFrom
			stringWriter
		}{rw, rw, rw, rw}
	case 0xc:
		return struct {
			http.ResponseWriter
			http.CloseNotifier
			http.Hijacker
		}{rw, rw, rw}
	case 0xd:
		return struct {
			http.ResponseWriter
			http.CloseNotifier
			http.Hijacker
			stringWriter
		}{rw, rw, rw, rw}
	case 0xe:
		return struct {
			http.ResponseWriter
			http.CloseNotifier
			http.Hijacker
			io.ReaderFrom
		}{rw, rw, rw, rw}
	case 0xf:
		return struct {
			http.ResponseWriter
			http.CloseNotifier
			http.Hijacker
			io.ReaderFrom
			stringWriter
		}{rw, rw, rw, rw, rw}
	case 0x10:
		return struct {
			http.ResponseWriter
			http.Flusher
		}{rw, rw}
	case 0x11:
		return struct {
			http.ResponseWriter
			http.Flusher
			stringWriter
		}{rw, rw, rw}
	case 0x12:
		return struct {
			http.ResponseWriter
			http.Flusher
			io.ReaderFrom
		}{rw, rw, rw}
	case 0x13:
		return struct {
			http.ResponseWriter
			http.Flusher
			io.ReaderFrom
			stringWriter
		}{rw, rw, rw, rw}
	case 0x14:
		return struct {
			http.ResponseWriter
			http.Flusher
			http.Hijacker
		}{rw, rw, rw}
	case 0x15:
		return struct {
			http.ResponseWriter
			http.Flusher
			http.Hijacker
			stringWriter
		}{rw, rw, rw, rw}
	case 0x16:
		return struct {
			http.ResponseWriter
			http.Flusher
			http.Hijacker
			io.ReaderFrom
		}{rw, rw, rw, rw}
	case 0x17:
		return struct {
			http.ResponseWriter
			http.Flusher
			http.Hijacker
			io.ReaderFrom
			stringWriter
		}{rw, rw, rw, rw, rw}
	case 0x18:
		return struct {
			http.ResponseWriter
			http.Flusher
			http.CloseNotifier
		}{rw, rw, rw}
	case 0x19:
		return struct {
			http.ResponseWriter
			http.Flusher
			http.CloseNotifier
			stringWriter
		}{rw, rw, rw, rw}
	case 0x1a:
		return struct {
			http.ResponseWriter
			http.Flusher
			http.CloseNotifier
			io.ReaderFrom
		}{rw, rw, rw, rw}
	case 0x1b:
		return struct {
			http.ResponseWriter
			http.Flusher
			http.CloseNotifier
			io.ReaderFrom
			stringWriter
		}{rw, rw, rw, rw, rw}
	case 0x1c:
		return struct {
			http.ResponseWriter
			http.Flusher
			http.CloseNotifier
			http.Hijacker
		}{rw, rw, rw, rw}
	case 0x1d:
		return struct {
			http.ResponseWriter
			http.Flusher
			http.CloseNotifier
			http.Hijacker
			stringWriter
		}{rw, rw, rw, rw, rw}
	case 0x1e:
		return struct {
			http.ResponseWriter
			http.Flusher
			http.CloseNotifier
			http.Hijacker
			io.ReaderFrom
		}{rw, rw, rw, rw, rw}
	case 0x1f:
		return struct {
			http.ResponseWriter
			http.Flusher
			http.CloseNotifier
			http.Hijacker
			io.ReaderFrom
			stringWriter
		}{rw, rw, rw, rw, rw, rw}
	}
	panic("unreachable")
}
