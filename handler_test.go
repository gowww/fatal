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

	if w.Code != http.StatusInternalServerError {
		t.Fail()
	}
}

func TestRecoverHandler(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	log.SetFlags(0)
	defer log.SetOutput(os.Stderr)

	status := http.StatusServiceUnavailable
	err := "unknown error"

	h := HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		panic(err)
	}, &Options{
		RecoverHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(status)
			err := Error(r)
			switch err.(type) {
			case string:
				w.Write([]byte(err.(string)))
			default:
				t.Errorf("error type: want string, got %T", err)
			}
		}),
	})
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != status {
		t.Errorf("status code: want %v, got %v", status, w.Code)
	}
	if w.Body.String() != err {
		t.Errorf("body: want %q, got %q", err, w.Body.String())
	}
}
