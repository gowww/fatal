package fatal

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestSimple(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	log.SetFlags(0)
	defer log.SetOutput(os.Stderr)

	h := HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("error")
	}, nil)
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fail()
	}
}

func TestRecoverHandler(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	log.SetFlags(0)
	defer log.SetOutput(os.Stderr)

	status := http.StatusInternalServerError
	body := http.StatusText(status)

	h := HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("error")
	}, &Options{
		RecoverHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(status)
			w.Write([]byte(body))
		}),
	})
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != status {
		t.Errorf("status code: want %v, got %v", status, w.Code)
	}
	if w.Body.String() != body {
		t.Errorf("body: want %q, got %q", body, w.Body.String())
	}
}
