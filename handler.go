// Package fatal provides a handler that recovers from panics.
package fatal

import (
	"bufio"
	"context"
	"log"
	"net"
	"net/http"
	"runtime"
)

type contextKey int

// Context keys
const (
	contextKeyError contextKey = iota
)

// A handler provides a handler that recovers from panics.
type handler struct {
	options *Options
	next    http.Handler
}

// Options provides the handler options.
type Options struct {
	RecoverHandler http.Handler // RecoverHandler, if provided, is called when recovering.
}

// Handle returns a Handler wrapping another http.Handler.
func Handle(h http.Handler, o *Options) http.Handler {
	return &handler{o, h}
}

// HandleFunc returns a Handler wrapping an http.HandlerFunc.
func HandleFunc(f http.HandlerFunc, o *Options) http.Handler {
	return Handle(f, o)
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil && err != http.ErrAbortHandler {
			buf := make([]byte, 64<<10)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("%v\n%s", err, buf)
			if !w.(*fatalWriter).written {
				if h.options != nil && h.options.RecoverHandler != nil {
					w.Header().Del("Content-Encoding")
					w.Header().Del("Content-Length")
					w.Header().Del("Content-Type")
					h.options.RecoverHandler.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), contextKeyError, err)))
				} else {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}
		}
	}()

	w = &fatalWriter{ResponseWriter: w}
	h.next.ServeHTTP(w, r)
}

// fatalWriter allows to tell if the response header has been written.
type fatalWriter struct {
	http.ResponseWriter
	written bool
}

func (fw *fatalWriter) WriteHeader(status int) {
	fw.written = true
	fw.ResponseWriter.WriteHeader(status)
}

func (fw *fatalWriter) Write(b []byte) (int, error) {
	fw.written = true
	return fw.ResponseWriter.Write(b)
}

// CloseNotify implements the http.CloseNotifier interface.
// No channel is returned if CloseNotify is not implemented by an upstream response writer.
func (fw *fatalWriter) CloseNotify() <-chan bool {
	n, ok := fw.ResponseWriter.(http.CloseNotifier)
	if !ok {
		return nil
	}
	return n.CloseNotify()
}

// Flush implements the http.Flusher interface.
// Nothing is done if Flush is not implemented by an upstream response writer.
func (fw *fatalWriter) Flush() {
	f, ok := fw.ResponseWriter.(http.Flusher)
	if ok {
		f.Flush()
	}
}

// Hijack implements the http.Hijacker interface.
// Error http.ErrNotSupported is returned if Hijack is not implemented by an upstream response writer.
func (fw *fatalWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := fw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, http.ErrNotSupported
	}
	return h.Hijack()
}

// Push implements the http.Pusher interface.
// http.ErrNotSupported is returned if Push is not implemented by an upstream response writer or not supported by the client.
func (fw *fatalWriter) Push(target string, opts *http.PushOptions) error {
	p, ok := fw.ResponseWriter.(http.Pusher)
	if !ok {
		return http.ErrNotSupported
	}
	return p.Push(target, opts)
}

// Error returns the error value stored in request's context during recovering.
func Error(r *http.Request) interface{} {
	return r.Context().Value(contextKeyError)
}
