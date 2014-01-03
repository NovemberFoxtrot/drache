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

	for server := range layout[os.Args[2]].Servers {
		theScripts := layout[os.Args[2]].Servers[server]

		fmt.Println(server)

		for _, theScript := range theScripts {
			fmt.Printf("  %s: ", theScript)
			script := &scripts.Script{Command: os.Args[1], Environment: os.Args[2], Status: 0}

			fmt.Println(script)

			//err := script.ParseLayout()

			//if err != nil {
			//	panic(err)
			//}

			//script.Run()
		}
	}

	//os.Exit(script.Status)
}
