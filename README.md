Serves a directory over HTTP
============================

I was extremely annoyed at *python -m SimpleHTTPServer* (lack of) speed so I
wrote one.

This project depends only on stdlib on purpose.


Installation
------------

    go get github.com/maruel/serve-dir


Usage
-----

Serve the current directory:

    serve-dir

Help with the command line arguments available:

    serve-dir -help


Example output
--------------
    11:15:52.282045 Serving /home/my_account/src on port 8010
    11:15:53.916813 192.168.1.2:2092 - 304      0b  GET /src/
    11:15:54.010258 192.168.1.2:2092 - 404     19b  GET /favicon.ico
    11:16:08.770496 192.168.1.2:2094 - 200   8877b  GET /src/foo.json
