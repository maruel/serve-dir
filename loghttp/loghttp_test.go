// Copyright 2017 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Example and tests for HTTP request logging and response forwarding.

package loghttp_test

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/maruel/serve-dir/loghttp"
)

func ExampleHandler() {
	// Serves the current directory over HTTP and logs all requests.
	log.SetFlags(log.Lmicroseconds)
	s := &http.Server{
		Addr:           ":6060",
		Handler:        &loghttp.Handler{Handler: http.FileServer(http.Dir("."))},
		ReadTimeout:    10. * time.Second,
		WriteTimeout:   24 * 60 * 60 * time.Second,
		MaxHeaderBytes: 256 * 1024 * 1024 * 1024,
	}
	log.Fatal(s.ListenAndServe())
}

func TestServeHTTP(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/foo", &bytes.Buffer{})
	h := loghttp.Handler{Handler: &dummy{}}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	r, err := io.ReadAll(w.Result().Body)
	if err != nil {
		t.Fatal(err)
	}
	if s := string(r); s != "hello" {
		t.Fatalf("%q != \"hello\"", s)
	}
}

func TestHandlerWithoutHijacker(t *testing.T) {
	t.Parallel()
	h := loghttp.Handler{Handler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if _, ok := w.(http.Hijacker); ok {
			t.Error("wrapper advertises hijacking when the underlying writer cannot hijack")
		}
	})}
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil))
}

type dummy struct {
}

func (d *dummy) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	_, _ = io.WriteString(w, "hello")
}

func init() {
	log.SetOutput(io.Discard)
}
