package main

import (
  "encoding/json"
  "fmt"
  "log"
  "net/http"
  "time"
  "strconv"
  "github.com/justinas/alice"
  "github.com/gorilla/mux"
)

func simpleLog(w http.ResponseWriter, r *http.Request) {
  // Shame on you! Make custom wrapper to HTTPHandler or something similar...
  // But if that's working then ok for now :)
  timestamp, _ := strconv.Atoi(w.Header().Get("Request-Start"))
  delta := time.Now().Unix() - int64(timestamp)
  log.Println(r.RemoteAddr, r.Method, r.URL, 200, delta)
}

func setHeaders(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // it's a very lame way of storing timestamp. There should be some other way to do this. For now it's not important.
    w.Header().Set("Request-Start", fmt.Sprintf("%v", time.Now().Unix()))
    w.Header().Set("Content-Type", "application/json")
    next.ServeHTTP(w, r)
  })
}

// JSON Response

type Response map[string]interface{}

func (r Response) String() (s string) {
  b, err := json.Marshal(r)
  if err != nil {
    s = ""
    return
  }
  s = string(b)
  return
}

func send(w http.ResponseWriter, r *http.Request, json Response) {
  w.Header().Set("Content-Lenght", fmt.Sprintf("%v",len(json)))
  fmt.Fprint(w, json)
  simpleLog(w, r)
}

//Controllers

// root
func rootController(w http.ResponseWriter, r *http.Request) {
  send(w, r, Response{"success": true, "message": "echodb http server is running!"})
}

// list all of collections
func collectionsController(w http.ResponseWriter, r *http.Request) {
  send(w, r, Response{"success": true, "message": "there should be a list of collections"})
}

// get a collection by name
func collectionController(w http.ResponseWriter, r *http.Request) {
  params := mux.Vars(r)
  send(w, r, Response{"success": true, "message": "collection: " + params["name"]})
}

// create a collection
func newCollectionController(w http.ResponseWriter, r *http.Request) {
  send(w, r, Response{"success": true, "message": "you have created a collection"})
}

// delete a collection
func deleteCollectionController(w http.ResponseWriter, r *http.Request) {
  send(w, r, Response{"success": true, "message": "you have deleted the collection"})
}

// list documents
func documentsController(w http.ResponseWriter, r *http.Request) {
  params := mux.Vars(r)

  send(w, r, Response{"success": true, "message": "list of documents in: " + params["name"]})
}

// read document
func documentController(w http.ResponseWriter, r *http.Request) {
  params := mux.Vars(r)

  send(w, r, Response{"success": true, "message": "a document " + params["id"] + " in: " + params["name"]})
}

// read document
func newDocumentController(w http.ResponseWriter, r *http.Request) {
  params := mux.Vars(r)

  send(w, r, Response{"success": true, "message": "new document" + " in: " + params["name"]})
}

// update document
func updateDocumentController(w http.ResponseWriter, r *http.Request) {
  params := mux.Vars(r)

  send(w, r, Response{"success": true, "message": "update document " + params["id"] + " in: " + params["name"]})
}

// delete document
func deleteDocumentController(w http.ResponseWriter, r *http.Request) {
  params := mux.Vars(r)

  send(w, r, Response{"success": true, "message": "delete document " + params["id"] + " in: " + params["name"]})
}

// ROUTER
func router() {
  stdChain := alice.New(setHeaders)

  r := mux.NewRouter()

  r.Handle("/", stdChain.Then(http.HandlerFunc(rootController)))

  // collection routers
  r.Handle("/colls", stdChain.Then(http.HandlerFunc(collectionsController)))
  r.Handle("/colls/{name}", stdChain.Then(http.HandlerFunc(collectionController)))
  r.Handle("/colls", stdChain.Then(http.HandlerFunc(newCollectionController))).Methods("POST")
  r.Handle("/colls/{name}", stdChain.Then(http.HandlerFunc(deleteCollectionController))).Methods("DELETE")

  // document routers
  r.Handle("/colls/{name}/docs", stdChain.Then(http.HandlerFunc(documentsController)))
  r.Handle("/colls/{name}/docs/{id}", stdChain.Then(http.HandlerFunc(documentController)))
  r.Handle("/colls/{name}/docs", stdChain.Then(http.HandlerFunc(newDocumentController))).Methods("POST")
  r.Handle("/colls/{name}/docs/{id}", stdChain.Then(http.HandlerFunc(updateDocumentController))).Methods("PUT")
  r.Handle("/colls/{name}/docs/{id}", stdChain.Then(http.HandlerFunc(deleteDocumentController))).Methods("DELETE")

  http.Handle("/", r)
  return
}

// main function
func main() {
  router()
  port := ":8001"
  log.Println("[HTTP Server]", port)
  http.ListenAndServe(port, nil)
}
