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

type Template struct {
	attributes map[string]interface{}
	source     string
}

func (t *Template) render() string {
	// ERB.new(t.source).result(binding)
	return t.source
}

type Layout struct {
	Attributes map[string]interface{} `json:"attributes"`
	Servers    map[string][]string    `json:"servers"`
}

func layout(directory string) map[string]Layout {
	input, err := ioutil.ReadFile(path.Join(directory, "layout.json"))

	if err != nil {
		fmt.Println(err)
	}

	var e map[string]Layout

	json.Unmarshal(input, &e)

	return e
}

func run(directory, server, recipe, command string, attributes map[string]interface{}) (string, int) {
	template_path := path.Join(directory, "recipe", recipe, command)

	if _, err := os.Stat(template_path); os.IsNotExist(err) {
		fmt.Print("\033[01;33mMISSING\033[00m ")
		return "unable to locate: " + template_path, 1
	}

	source, err := ioutil.ReadFile(template_path)

	if err != nil {
		return "unable to read file: " + template_path, 1
	}

	template := &Template{source: string(source), attributes: attributes}

	out, status := ssh(server, template.render())

	return out, status
}

func ssh(server, script string) (string, int) {
	cmd := exec.Command(script[:len(script)-1]) // TODO

	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Print(err)
		return string(output), 1
	}

	// out, status = Open3.capture2e("ssh -T -F #{path("ssh_config")} #{server}", :stdin_data => script)

	return string(output), 0
}

func main() {
	var command = flag.String("c", "", "command")
	var directory = flag.String("d", ".", "directory")
	var environment = flag.String("e", "", "environment")
	var verbose = flag.Bool("v", false, "verbose mode")

	flag.Parse()

	e := layout(*directory)

	layout := e[*environment]

	exit_status := 0

	for v, _ := range layout.Servers {
		recipes := layout.Servers[v]

		fmt.Println(v)

		for _, recipe := range recipes {
			fmt.Printf("  %s: ", recipe)

			filename := path.Join(*directory, "recipe", recipe)

			if _, err := os.Stat(filename); os.IsNotExist(err) {
				fmt.Printf("unable to locate: %s\n", filename)
				os.Exit(1)
			}

			stdout, status := run(*directory, v, recipe, *command, layout.Attributes)

			switch status {
			case -1: // nil is better? negative error codes?
				fmt.Print("?")
			case 0:
				fmt.Print("\033[01;32mOK\033[00m\n")
			default:
				fmt.Print("\033[01;31mERROR\033[00m\n")
				exit_status = 1
				break
			}

			if *verbose && len(stdout) > 0 {
				fmt.Fprintf(os.Stderr, " %s\n", stdout)
			}
		}
	}

	fmt.Println("\033[01;32mDONE\033[00m")
	os.Exit(exit_status)
}
