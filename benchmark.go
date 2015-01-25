package main

import (
	"./db"
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

	for i := 0; i < 1000000; i++ {
		_, err := books.Insert(map[string]interface{}{
			"name":   "An introduction to programming in Go",
			"author": "Caleb Doxsey"})
		if err != nil {
			panic(err)
		}
	}

	echodb.Close()

}
