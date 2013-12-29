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

type Recipe struct {
	attributes map[string]interface{}
	source     string
}

func (r *Recipe) render() string {
	return r.source
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

	err = json.Unmarshal(input, &b.layout)

	if err != nil {
		fmt.Println(err)
	}
}

func (b *Book) exec(server, recipe string) (string, int) {
	template_path := path.Join(".", "recipe", recipe, b.command)

	if _, err := os.Stat(template_path); os.IsNotExist(err) {
		fmt.Print("\033[01;33mMISSING\033[00m ")
		return "unable to locate: " + template_path, 1
	}

	/*
		source, err := ioutil.ReadFile(template_path)

		if err != nil {
			return "unable to read file: " + template_path, 1
		}
	*/
	// template := &Recipe{source: string(source), attributes: attributes}

	// out, status := ssh(server, template.render())
	out, status := ssh(server, template_path)

	return out, status
}

func ssh(server, script string) (string, int) {
	cmd := exec.Command(script)
	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Print(err)
		return string(output), 1
	}

	// out, status = Open3.capture2e("ssh -T -F #{path("ssh_config")} #{server}", :stdin_data => script)
	return string(output), 0
}

func (b *Book) run() {
	for server := range b.layout[b.environment].Servers {
		recipes := b.layout[b.environment].Servers[server]

		fmt.Println(server)

		for _, recipe := range recipes {
			fmt.Printf("  %s: ", recipe)

			stdout, status := b.exec(server, recipe)

			if status != 0 && len(stdout) > 0 {
				fmt.Fprintf(os.Stderr, " %s\n", stdout)
			}

			if status != 0 {
				fmt.Print("\033[01;31mERROR\033[00m\n")
				b.status = 1
				break
			}

			fmt.Print("\033[01;32mOK\033[00m\n")
		}
	}
}

func main() {
	book := &Book{command: os.Args[1], environment: os.Args[2], status: 0}

	book.ParseLayout()

	book.run()

	os.Exit(book.status)
}
