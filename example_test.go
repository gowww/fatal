package fatal_test

import (
	"github.com/gowww/fatal"
	"net/http"
)

func Example() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		panic("error")
	})

	http.ListenAndServe(":8080", fatal.Handle(mux, &fatal.Options{
		RecoverHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}),
	}))
}

func ExampleHandle() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		panic("error")
	})

	http.ListenAndServe(":8080", fatal.Handle(mux, &fatal.Options{
		RecoverHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}),
	}))
}

func ExampleHandleFunc() {
	http.Handle("/", fatal.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("error")
	}, &fatal.Options{
		RecoverHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}),
	}))

	http.ListenAndServe(":8080", nil)
}
