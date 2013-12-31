package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
)

type Book struct {
	command     string
	environment string
	layout      map[string]Layout
	status      int
}

type Layout struct {
	Attributes map[string]interface{} `json:"attributes"`
	Servers    map[string][]string    `json:"servers"`
}

func (b *Book) ParseLayout() {
	input, err := ioutil.ReadFile(path.Join(".", "layout.json"))

	if err != nil {
		fmt.Println(err)
	}

	if err = json.Unmarshal(input, &b.layout); err != nil {
		fmt.Println(err)
	}
}

func (b *Book) exec(server, recipe string) (string, int) {
	scriptPath := path.Join(".", "recipe", recipe, b.command)

	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return "\033[01;33mMISSING\033[00m unable to locate: " + scriptPath, 0
	}

	if source, err := ioutil.ReadFile(scriptPath); err != nil {
		return "unable to read file: " + scriptPath, 1
	} else {
		out, status := ssh(server, string(source))
		return out, status
	}
}

func ssh(server, script string) (string, int) {
	cmd := exec.Command("ssh", "-T", server, script)

	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Print(err)
		return string(output), 1
	}

	return string(output), 0
}

func (b *Book) run() {
	for server := range b.layout[b.environment].Servers {
		scripts := b.layout[b.environment].Servers[server]

		fmt.Println(server)

		for _, script := range scripts {
			fmt.Printf("  %s: ", script)

			stdout, status := b.exec(server, script)

			if status != 0 {
				fmt.Fprintf(os.Stderr, "\033[01;31mERROR\033[00m\n %s\n", stdout)
				b.status = 1
				break
			}

			fmt.Print("\033[01;32mOK\033[00m\n")
		}
	}
}

func main() {
	book := &Book{command: os.Args[2], environment: os.Args[1], status: 0}
	book.ParseLayout()
	book.run()
	os.Exit(book.status)
}
