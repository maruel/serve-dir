/* Copyright 2012 Marc-Antoine Ruel. Licensed under the Apache License, Version
2.0 (the "License"); you may not use this file except in compliance with the
License.  You may obtain a copy of the License at
http://www.apache.org/licenses/LICENSE-2.0. Unless required by applicable law or
agreed to in writing, software distributed under the License is distributed on
an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
or implied. See the License for the specific language governing permissions and
limitations under the License. */

/* Serves a directory over HTTP and logs the request to stderr. */
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var port = flag.Int("port", 8010, "port number")
var root_dir = flag.String("root", getWd(), "root directory")
var timeout = flag.Int("timeout", 24*60*60, "write timeout in seconds; default 24h")
var max_size = flag.Int("max_size", 256*1024*1024*1024, "max transfer size; default 256gb")

func getWd() string {
	wd, _ := os.Getwd()
	return wd
}

type loggingHandler struct {
	handler http.Handler
}

type loggingResponseWriter struct {
	http.ResponseWriter
	length int
	status int
}

func (l *loggingResponseWriter) Write(data []byte) (size int, err error) {
	size, err = l.ResponseWriter.Write(data)
	l.length += size
	return
}

func (l *loggingResponseWriter) WriteHeader(status int) {
	l.ResponseWriter.WriteHeader(status)
	l.status = status
}

// Logs each HTTP request.
func (l loggingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	l_w := &loggingResponseWriter{ResponseWriter: w}
	l.handler.ServeHTTP(l_w, r)
	log.Printf("%s - %3d %6db %4s %s\n",
		r.RemoteAddr,
		l_w.status,
		l_w.length,
		r.Method,
		r.RequestURI)
}

func main() {
	log.SetFlags(log.Lmicroseconds)
	flag.Parse()

	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: loggingHandler{http.FileServer(http.Dir(*root_dir))},
		// read timeout is always 10s, since it should be GETs only.
		ReadTimeout:    10. * time.Second,
		WriteTimeout:   time.Duration(*timeout) * time.Second,
		MaxHeaderBytes: *max_size,
	}
	log.Printf("Serving %s on port %d", *root_dir, *port)
	log.Fatal(s.ListenAndServe())
}
