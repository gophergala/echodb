package main

import (
  "encoding/json"
  "fmt"
  "log"
  "net/http"
  "github.com/justinas/alice"
  "github.com/gorilla/mux"
)

/*
router

COLLECTIONS
get /db/collections/
get /db/collections/:collection
post /db/collections/
delete /db/collections/:collection

DOCUMENTS
get /db/collections/:collection/docs/
get db/colls/:collection/docs/:id
post /db/colls/:collection/docs
put /db/colls/:collection/docs/:id
delete /db/colls/:collection/docs/:id
*/

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

func rootController(w http.ResponseWriter, r *http.Request) {
  // params := mux.Vars(r)
  w.Header().Set("Content-Type", "application/json")
  fmt.Fprint(w, Response{"success": true, "message": "echodb http server is running!"})
}

func router() {
  stdChain := alice.New()

  r := mux.NewRouter()

  r.Handle("/", stdChain.Then(http.HandlerFunc(rootController)))
  // r.Handle("/colls", stdChain.Then(http.HandlerFunc(collectionsController)))
  // r.Handle("/colls/:name", stdChain.Then(http.HandlerFunc(collectionController)))
  // r.Handle("/colls/:name", stdChain.Then(http.HandlerFunc(newCollectionController))).Methods("POST")
  // r.Handle("/colls/:name", stdChain.Then(http.HandlerFunc(deleteCollectionController))).Methods("DELETE")

  http.Handle("/", r)
  return
}

func main() {
  router()

  http.ListenAndServe(":8001", nil)
  log.Println("Server's up and running!")

}
