package main

import (
	"fmt"
	"os"

	"github.com/NovemberFoxtrot/remote/layouts"
	"github.com/NovemberFoxtrot/remote/scripts"
)

func main() {

	theCommand := os.Args[1]
	theDirectory := os.Args[3]
	theEnvironment := os.Args[2]

	input, err := layouts.Read(theDirectory)

	if err != nil {
		panic(err)
	}

	var layout layouts.Layout

	err = layout.Parse(input)

	if err != nil {
		panic(err)
	}

	fmt.Println(layout)

	status := 0

	for server := range layout[theEnvironment].Servers {
		theScripts := layout[theEnvironment].Servers[server]

		fmt.Println(server)

		for _, theScript := range theScripts {
			fmt.Printf("  %s: ", theScript)

			script := &scripts.Script{
				Command:     theCommand,
				Directory:   theDirectory,
				Environment: theEnvironment,
				Name:        theScript,
				Server:      server,
				Status:      0,
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
