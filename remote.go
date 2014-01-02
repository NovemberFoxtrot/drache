package main

import (
	"os"

	"github.com/NovemberFoxtrot/remote/scripts"
)

func main() {
	script := &scripts.Script{Command: os.Args[2], Environment: os.Args[1], Status: 0}

	err := script.ParseLayout()

	if err != nil {
		panic(err)
	}

	script.Run()

	os.Exit(script.Status)
}
