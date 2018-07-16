// Copyright 2017 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package loghttp exposes a http.Handler that logs requests via log.Printf().
package loghttp

import (
	"bufio"
	"log"
	"net"
	"net/http"
	"time"
)

// Handler wraps a http.Handler and logs the request and response.
//
// It handles Hijack() for websocket support.
type Handler struct {
	http.Handler

	unused struct{}
}

// ServeHTTP implements http.Handler.
//
// It logs the handler's HTTP status code and the response size.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	received := time.Now()
	rwh := responseWriteHijacker{responseWriter: responseWriter{ResponseWriter: w}}
	w = &rwh
	// Not all ResponseWriter implement Hijack, so query its support upfront.
	if _, ok := w.(http.Hijacker); !ok {
		w = &rwh.responseWriter
	}
	defer func() {
		defaultOnDone(r, received, rwh.status, rwh.length, rwh.hijacked)
	}()
	h.Handler.ServeHTTP(w, r)
}

// Private.

func defaultOnDone(r *http.Request, received time.Time, status, length int, hijacked bool) {
	m := r.Method
	if hijacked {
		m = "HIJACKED"
	}
	log.Printf("%s - %3d %6db %-4s %6s %s", r.RemoteAddr, status, length, m, roundDuration(time.Since(received)), r.RequestURI)
}

type responseWriter struct {
	http.ResponseWriter
	length int
	status int
}

func (r *responseWriter) Write(data []byte) (size int, err error) {
	if r.status == 0 {
		r.status = 200
	}
	size, err = r.ResponseWriter.Write(data)
	r.length += size
	return
}

func (r *responseWriter) WriteHeader(status int) {
	r.ResponseWriter.WriteHeader(status)
	r.status = status
}

type responseWriteHijacker struct {
	responseWriter
	hijacked bool
}

// Hijack is needed for websocket.
//
// TODO(maruel): For now the write length is lost, it would require querying
// the ReadWriter and overriding Conn.Write().
func (r *responseWriteHijacker) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	r.hijacked = true
	return r.ResponseWriter.(http.Hijacker).Hijack()
}

//

// roundDuration returns time rounded to 4 digits.
func roundDuration(d time.Duration) time.Duration {
	if l := log10(int64(d)); l > 3 {
		m := time.Duration(1)
		for i := uint(3); i < l; i++ {
			m *= 10
		}
		d = (d + (m / 2)) / m * m
	}
	return d
}

// log10 is a cheap way to find the most significant digit.
func log10(i int64) uint {
	switch {
	case i < 10:
		return 0
	case i < 100:
		return 1
	case i < 1000:
		return 2
	case i < 10000:
		return 3
	case i < 100000:
		return 4
	case i < 1000000:
		return 5
	case i < 10000000:
		return 6
	case i < 100000000:
		return 7
	case i < 1000000000:
		return 8
	case i < 10000000000:
		return 9
	case i < 100000000000:
		return 10
	case i < 1000000000000:
		return 11
	case i < 10000000000000:
		return 12
	case i < 100000000000000:
		return 13
	case i < 1000000000000000:
		return 14
	case i < 10000000000000000:
		return 15
	default:
		return 16
	}
}
