EchoDB
===================
Fault-tolerrent "data-on-wire" NoSQL datastore written in 48hrs.

* Fault-tolerrant NoSQL [done]
* MMAP based datastore (mostly based on gommap and tiedot wrapper)
  [done]
* Hashtable based indexer (based on tiedot implementation) [done]
* Simple HTTP API to manage collections [done]
* Data on wire (using websocket) [almost done]
* and yes, it's written in 48hrs during
  [GopherGala](http://gophergala.com/) 2015

Install
===================
```
mkdir echodb
cd echodb
export GOPATH=`pwd`
go get github.com/gophergala/echodb
cp -r src/github.com/gophergala/echodb/todoapp bin/todoapp
./bin/echodb
```

There is a sample todo app at http://localhost:8001/client [try it in
two browser sessions]

HTTP server runs at http://localhost:8001 please see
[server.go](dbhttp/server.go)

Dependencies
======================
```
go get github.com/justinas/alice
go get github.com/gorilla/mux
go get github.com/gorilla/websocket
```

Current status
==================
Highly experimental

