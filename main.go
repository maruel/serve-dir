// Copyright 2017 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// serve-dir serves a directory over HTTP and logs the request to stderr.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/maruel/serve-dir/loghttp"
)

func getWd() string {
	wd, _ := os.Getwd()
	return wd
}

func main() {
	// TODO(maruel): Change to -http, so it can bind to localhost.
	port := flag.Int("port", 8010, "port number")
	rootDir := flag.String("root", getWd(), "root directory")
	timeout := flag.Int("timeout", 24*60*60, "write timeout in seconds; default 24h")
	maxHdrSize := flag.Int("max_size", http.DefaultMaxHeaderBytes, "max header transfer size")

	log.SetFlags(log.Lmicroseconds)
	flag.Parse()

	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: &loghttp.Handler{Handler: http.FileServer(http.Dir(*rootDir))},
		// read timeout is always 10s, since it should be GETs only.
		ReadTimeout:    10. * time.Second,
		WriteTimeout:   time.Duration(*timeout) * time.Second,
		MaxHeaderBytes: *maxHdrSize,
	}
	log.Printf("Serving %s on port %d", *rootDir, *port)
	log.Fatal(s.ListenAndServe())
}
