package main

import (
	"fmt"
	"os"

	"github.com/NovemberFoxtrot/remote/layouts"
	"github.com/NovemberFoxtrot/remote/scripts"
)

func main() {
	input, err := layouts.Read(os.Args[3])

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

	for server := range layout[os.Args[2]].Servers {
		theScripts := layout[os.Args[2]].Servers[server]

		fmt.Println(server)

		for _, theScript := range theScripts {
			fmt.Printf("  %s: ", theScript)
			script := &scripts.Script{Command: os.Args[1], Environment: os.Args[2], Status: 0}

			stdout, status := script.Run(server, os.Args[1])

			if status != 0 {
				fmt.Fprintf(os.Stderr, "\033[01;31mERROR\033[00m\n %s\n", stdout)
				break
			}

			fmt.Print("\033[01;32mOK\033[00m\n")
		}
	}

	os.Exit(status)
}
