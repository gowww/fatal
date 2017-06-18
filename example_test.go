package fatal_test

import (
	"fmt"
	"net/http"

	"github.com/gowww/fatal"
)

func Example() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		panic("error")
	})

	http.ListenAndServe(":8080", fatal.Handle(mux, &fatal.Options{
		RecoverHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, fmt.Sprintf("error: %v", fatal.Error(r)), http.StatusInternalServerError)
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
			http.Error(w, fmt.Sprintf("error: %v", fatal.Error(r)), http.StatusInternalServerError)
		}),
	}))
}

func ExampleHandleFunc() {
	http.Handle("/", fatal.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("error")
	}, &fatal.Options{
		RecoverHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, fmt.Sprintf("error: %v", fatal.Error(r)), http.StatusInternalServerError)
		}),
	}))

	http.ListenAndServe(":8080", nil)
}
