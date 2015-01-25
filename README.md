EchoDB
===================
Fault-tolerrent "data-on-wire" NoSQL datastore written in 48hrs.

* Fault-tolerrant NoSQL
* MMAP based datastore (mostly based on gommap and tiedot wrapper)
* Hashtable based indexer (based on tiedot implementation)
* Simple HTTP API to manage collections
* Data on wire (using websocket)
* and yes, it's written in 48hrs during
  [GopherGala](http://gophergala.com/) 2015

Install
===================
```
git clone git@github.com:gophergala/echodb.git
cd echodb
export GOPATH=`pwd`


go get github.com/justinas/alice
go get github.com/gorilla/mux
go get github.com/gorilla/websocket
```

Current status
==================
Highly experimental

