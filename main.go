package main

import (
	"fmt"
	"os"

	"github.com/NovemberFoxtrot/remote/layouts"
	// "github.com/NovemberFoxtrot/remote/scripts"
)

func main() {
	input, err := layouts.Read(os.Args[1])

	if err != nil {
		panic(err)
	}

	var layout layouts.Layout

	err = layout.Parse(input)

	if err != nil {
		panic(err)
	}

	fmt.Println(layout)

	//for server := range b.layout[b.Environment].Servers {
	//scripts := b.layout[b.Environment].Servers[server]

	//fmt.Println(server)
	//for _, script := range scripts {
	// fmt.Printf("  %s: ", script)
	// }
	//}

	//script := &scripts.Script{Command: os.Args[2], Environment: os.Args[1], Status: 0}

	//err := script.ParseLayout()

	//if err != nil {
	//	panic(err)
	//}

	//script.Run()

	//os.Exit(script.Status)
}
