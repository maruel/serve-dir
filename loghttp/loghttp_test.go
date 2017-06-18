// Copyright 2017 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package loghttp

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func ExampleHandler() {
	// Serves the current directory over HTTP and logs all requests.
	log.SetFlags(log.Lmicroseconds)
	s := &http.Server{
		Addr:           ":6060",
		Handler:        &Handler{Handler: http.FileServer(http.Dir("."))},
		ReadTimeout:    10. * time.Second,
		WriteTimeout:   24 * 60 * 60 * time.Second,
		MaxHeaderBytes: 256 * 1024 * 1024 * 1024,
	}
	log.Fatal(s.ListenAndServe())
}

func TestServeHTTP(t *testing.T) {
	req := httptest.NewRequest("GET", "/foo", &bytes.Buffer{})
	h := Handler{Handler: &dummy{}}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	r, _ := ioutil.ReadAll(w.Result().Body)
	if s := string(r); s != "hello" {
		t.Fatalf("%q != \"hello\"", s)
	}
}

type dummy struct {
}

func (d *dummy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello")
}
