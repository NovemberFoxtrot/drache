package main

import (
	"os"

	"github.com/NovemberFoxtrot/drache/scripts"
)

func main() {
	book := &scripts.Book{Command: os.Args[2], Environment: os.Args[1], Status: 0}
	book.ParseLayout()

	if err != nil {
		panic(err)
	}

	book.Run()

	os.Exit(book.Status)
}
