package dbhttp

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"log"
	"net/http"
	"time"
	"../dbwebsocket"
	"text/template"
	"../db"
)

func simpleLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()

		log.Println(r.RemoteAddr, r.Method, r.URL, 200, t2.Sub(t1))
	})
}

func recoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				http.Error(w, http.StatusText(500), 500)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func setHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	fmt.Fprint(w, json)
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
 	col := echodb.Get(params["name"])

	send(w, r, Response{"success": true, "count": fmt.Sprintf("%v", col.Count())})
}

// create a collection
func newCollectionController(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	echodb.Create(params["name"])

	col := echodb.Get(params["name"])
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

func serveWs(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	dbwebsocket.ServeWs(params["name"], w, r)
}


var homeTempl = template.Must(template.ParseFiles("./dbhttp/index.html"))
func serveHome(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	homeTempl.Execute(w, r.Host)
}


// ROUTER
func router() {
	stdChain := alice.New(simpleLogger, recoverHandler, setHeaders)

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

	r.Handle("/ws/{name}", http.HandlerFunc(serveWs)).Methods("GET")
	r.Handle("/client", http.HandlerFunc(serveHome)).Methods("GET")

	http.Handle("/", r)
	return
}

var echodb *db.Database
// main function
func Start() {
	echodb, _ = db.OpenDatabase("/tmp/echodb")
	router()
	port := ":8001"
	log.Println("[HTTP Server]", port)
	http.ListenAndServe(port, nil)
}
