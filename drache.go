package main

import (
	"os"

	"github.com/NovemberFoxtrot/drache/scripts"
)


func main() {
	book := &Book{command: os.Args[2], environment: os.Args[1], status: 0}
	book.ParseLayout()
	book.run()
	os.Exit(book.status)
}
