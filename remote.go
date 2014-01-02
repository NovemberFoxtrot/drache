package main

import (
	"os"

	"github.com/NovemberFoxtrot/remote/remote"
)

func main() {
	remote := &remote.Script{Command: os.Args[2], Environment: os.Args[1], Status: 0}

	err := remote.ParseLayout()

	if err != nil {
		panic(err)
	}

	remote.Run()

	os.Exit(remote.Status)
}
