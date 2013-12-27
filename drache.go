package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
)

type Book struct {
	command     string
	directory   string
	environment string
	verbose     bool
	status      int
}

type Recipe struct {
	attributes map[string]interface{}
	source     string
}

func (r *Recipe) render() string {
	// ERB.new(r.source).result(binding)

	return r.source
}

type Layout struct {
	Attributes map[string]interface{} `json:"attributes"`
	Servers    map[string][]string    `json:"servers"`
}

func (b *Book) layout() map[string]Layout {
	input, err := ioutil.ReadFile(path.Join(b.directory, "layout.json"))

	if err != nil {
		fmt.Println(err)
	}

	var layout map[string]Layout

	json.Unmarshal(input, &layout)

	return layout
}

func run(directory, server, recipe, command string, attributes map[string]interface{}) (string, int) {
	template_path := path.Join(directory, "recipe", recipe, command)

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
	jsonstruct := b.layout()
	layout := jsonstruct[b.environment]

	// TODO abort if layout does not parse

	for server := range layout.Servers {
		recipes := layout.Servers[server]

		fmt.Println(server)

		for _, recipe := range recipes {
			fmt.Printf("  %s: ", recipe)

			stdout, status := run(b.directory, server, recipe, b.command, layout.Attributes)

			if b.verbose && len(stdout) > 0 {
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
	var command = flag.String("c", "", "command")
	var directory = flag.String("d", ".", "directory")
	var environment = flag.String("e", "", "environment")
	var verbose = flag.Bool("v", false, "verbose mode")

	flag.Parse()

	book := &Book{command: *command, directory: *directory, environment: *environment, verbose: *verbose, status: 0}

	book.run()

	fmt.Println("\033[01;32mDONE\033[00m")
	os.Exit(book.status)
}
