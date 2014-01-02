package main

import (
	"os"

	"github.com/NovemberFoxtrot/remote/layouts"
	"github.com/NovemberFoxtrot/remote/scripts"
)

func main() {
	input, err := layouts.ReadLayout()

	if err != nil {
		panic(err)
	}

	layout, err := layouts.ParseLayout(input)

	//for server := range b.layout[b.Environment].Servers {
		//scripts := b.layout[b.Environment].Servers[server]

		//fmt.Println(server)
		//for _, script := range scripts {
	// fmt.Printf("  %s: ", script)
		}
	//}

	script := &scripts.Script{Command: os.Args[2], Environment: os.Args[1], Status: 0}

	err := script.ParseLayout()

	if err != nil {
		panic(err)
	}

	script.Run()

	os.Exit(script.Status)
}
