// +build go1.8

package httplogger

import (
	"net/http"
)

func (rw *responseWriter) Push(target string, opts *http.PushOptions) error {
	if p, ok := rw.rw.(http.Pusher); ok {
		return p.Push(target, opts)
	}
	return http.ErrNotSupported
}
