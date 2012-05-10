A basic HTTP directory server with logging
==========================================

It was simply faster for me to code it that find it one on the interwebs.
I was extremely annoyed at python -m SimpleHTTPServer (lack of) speed.

Run ./serve-dir -help for help.

Example output
--------------
    11:15:52.282045 Serving /home/my_account/src on port 8010
    11:15:53.916813 192.168.1.2:2092 - 304      0b  GET /src/
    11:15:54.010258 192.168.1.2:2092 - 404     19b  GET /favicon.ico
    11:16:08.770496 192.168.1.2:2094 - 200   8877b  GET /src/foo.json
