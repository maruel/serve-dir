# Serves a directory over HTTP

I was extremely annoyed at *python -m SimpleHTTPServer* (lack of) speed so I
wrote one.

This project depends only on stdlib on purpose.


## Installation

    go install github.com/maruel/serve-dir@latest


## Usage

Serve the current directory:

    serve-dir

Help with the command line arguments available:

    serve-dir -help


## Example output

    11:15:52.282045 Serving /home/my_account/src on port 8010
    11:15:53.916813 192.168.1.2:2092 - 304      0b  GET /src/
    11:15:54.010258 192.168.1.2:2092 - 404     19b  GET /favicon.ico
    11:16:08.770496 192.168.1.2:2094 - 200   8877b  GET /src/foo.json


# Logging library

The [http.Handler](https://pkg.go.dev/net/http#Handler) logging code in
`serve-dir` is usable as a library as `github.com/maruel/serve-dir/loghttp` via
[loghttp.Handler](https://pkg.go.dev/github.com/maruel/serve-dir/loghttp#Handler).

[![Go
Reference](https://pkg.go.dev/badge/github.com/maruel/serve-dir/loghttp.svg)](https://pkg.go.dev/github.com/maruel/serve-dir/loghttp)

Example:

```go
// Serves the current directory over HTTP and logs all requests.
log.SetFlags(log.Lmicroseconds)
s := &http.Server{
    Addr:           ":6060",
    Handler:        &loghttp.Handler{Handler: http.FileServer(http.Dir("."))},
}
log.Fatal(s.ListenAndServe())
```
