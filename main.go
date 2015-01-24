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
	} else {
		echodb.Create("books")
	}
	dbhttp.Start()
}
