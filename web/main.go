package main

import (
	"context"
	"errors"
	"log"
	"net/http"
)

func main() {
	srv := newServer()

	http.Handle("/gen", contextHandler(srv.GenerateStateHandler))
	http.Handle("/state", contextHandler(srv.StateHandler))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

type contextHandler func(ctx context.Context, w http.ResponseWriter, req *http.Request) error

func (hndl contextHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	err := hndl(ctx, w, req)
	if err == nil {
		return
	}

	var he httpError
	if errors.As(err, &he) {
		http.Error(w, err.Error(), he.StatusCode())
		return
	}

	log.Printf("%s: %v", req.RequestURI, err)
	http.Error(w, "Internal server error", http.StatusInternalServerError)
}

type httpError struct {
	msg  string
	code int
}

func (e httpError) Error() string {
	return e.msg
}

func (e httpError) StatusCode() int {
	return e.code
}
