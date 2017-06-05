# [![gowww](https://avatars.githubusercontent.com/u/18078923?s=20)](https://github.com/gowww) fatal [![GoDoc](https://godoc.org/github.com/gowww/fatal?status.svg)](https://godoc.org/github.com/gowww/fatal) [![Build](https://travis-ci.org/gowww/fatal.svg?branch=master)](https://travis-ci.org/gowww/fatal) [![Coverage](https://coveralls.io/repos/github/gowww/fatal/badge.svg?branch=master)](https://coveralls.io/github/gowww/fatal?branch=master) [![Go Report](https://goreportcard.com/badge/github.com/gowww/fatal)](https://goreportcard.com/report/github.com/gowww/fatal)

Package [fatal](https://godoc.org/github.com/gowww/fatal) provides a handler that recovers from panics.

## Example

```Go
mux := http.NewServeMux()

mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	panic("error")
})

http.ListenAndServe(":8080", fatal.Handle(mux, &fatal.Options{
	RecoverHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}),
}))
```
