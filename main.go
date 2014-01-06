package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/NovemberFoxtrot/remote/layouts"
	"github.com/NovemberFoxtrot/remote/scripts"
)

func main() {
	theCommand := os.Args[1]
	theDirectory := os.Args[3]
	theEnvironment := os.Args[2]
	theVerbose := false

	if len(os.Args) > 4 && strings.Contains(os.Args[4], "-v") {
		theVerbose = true
	}

	var layout layouts.Layout

	input, err := layouts.Read(theDirectory)

	if err != nil {
		panic(err)
	}

	err = layout.Parse(input)

	if err != nil {
		panic(err)
	}

	status := 0

	for server := range layout[theEnvironment].Servers {
		theScripts := layout[theEnvironment].Servers[server]

		fmt.Println(server)

		for _, theScript := range theScripts {
			fmt.Printf(" %s: ", theScript)

			script := &scripts.Script{
				Attributes:  layout[theEnvironment].Attributes,
				Command:     theCommand,
				Directory:   theDirectory,
				Name:        theScript,
				Server:      server,
				Status:      0,
				Verbose:     theVerbose,
			}

			script.Run()

			if script.Status != 0 {
				fmt.Fprintf(os.Stderr, "\033[01;31mERROR\033[00m\n %s\n", script.Output)
				status = script.Status
				break
			}

			fmt.Print("\033[01;32mOK\033[00m\n")
		}
	}

	os.Exit(status)
}
