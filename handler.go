// Package fatal provides a handler that recovers from panics.
package fatal

import (
	"log"
	"net/http"
	"runtime"
)

// A handler provides a clever gzip compressing handler.
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
			if h.options != nil && h.options.RecoverHandler != nil && !w.(*fatalWriter).written {
				w.Header().Del("Content-Type")
				h.options.RecoverHandler.ServeHTTP(w, r)
			}
		}
	}()

	if h.options != nil && h.options.RecoverHandler != nil {
		w = &fatalWriter{ResponseWriter: w}
	}
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
