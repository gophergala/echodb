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

	doc, err := books.FindById(docId)
	if err != nil {
		panic(err)
	}

	fmt.Println("Document", docId, "is", doc)

	err = books.Update(docId, map[string]interface{}{
		"name":   "hack in go",
		"author": "you",
		"isbn":   "234238729837"})
	if err != nil {
		panic(err)
	}

	doc, err = books.FindById(docId)
	if err != nil {
		panic(err)
	}

	fmt.Println("Document", docId, "is", doc)
	count := books.Count()
	fmt.Println("Documents", count)

	// Gracefully close database
	if err := echodb.Close(); err != nil {
		panic(err)
	}

	dbhttp.Start()
}
