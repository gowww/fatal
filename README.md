# [![gowww](https://avatars.githubusercontent.com/u/18078923?s=20)](https://github.com/gowww) fatal [![GoDoc](https://godoc.org/github.com/gowww/fatal?status.svg)](https://godoc.org/github.com/gowww/fatal) [![Build](https://travis-ci.org/gowww/fatal.svg?branch=master)](https://travis-ci.org/gowww/fatal) [![Coverage](https://coveralls.io/repos/github/gowww/fatal/badge.svg?branch=master)](https://coveralls.io/github/gowww/fatal?branch=master) [![Go Report](https://goreportcard.com/badge/github.com/gowww/fatal)](https://goreportcard.com/report/github.com/gowww/fatal)

Package [fatal](https://godoc.org/github.com/gowww/fatal) provides a handler that recovers from panics.

## Installing

1. Get package:

	```Shell
	go get -u github.com/gowww/fatal
	````

2. Import it in your code:

	```Go
	import "github.com/gowww/fatal"
	```

## Usage

To wrap an [http.Handler](https://golang.org/pkg/net/http/#Handler), use [Handle](https://godoc.org/github.com/gowww/fatal#Handle):

```Go
mux := http.NewServeMux()

mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	panic("error")
})

http.ListenAndServe(":8080", fatal.Handle(mux, nil))
````

To wrap an [http.HandlerFunc](https://golang.org/pkg/net/http/#HandlerFunc), use [HandleFunc](https://godoc.org/github.com/gowww/fatal#HandleFunc):

```Go
http.Handle("/", fatal.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
	panic("error")
}, nil))

http.ListenAndServe(":8080", nil)
```

### Custom "error" handler

When a your code panics, a 500 error with an empty body is send by default.

But you can set your own "error" handler (and send an HTML page, for example):

```Go
http.ListenAndServe(":8080", fatal.Handle(mux, &fatal.Options{
	RecoverHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error page")
	}),
}))
```

Note that is this case, it's up to you to set the correct status code (normally 500) for the response.
