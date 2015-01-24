package main

import (
	"./db"
	"./dbhttp"
	"fmt"
)

func main() {
	echodb, err := db.OpenDatabase("/tmp/echodb")
	if err != nil {
		fmt.Printf("Error: %v", err)
		panic("exit")
	}

	echodb.Create("books")
	books := echodb.Get("books")

	docId, err := books.Insert(map[string]interface{}{
		"name":   "An introduction to programming in Go",
		"author": "Caleb Doxsey"})
	if err != nil {
		panic(err)
	}
	fmt.Println("DocumentID")
	fmt.Println(docId)

	dbhttp.Start()
}
