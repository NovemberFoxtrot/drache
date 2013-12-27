package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
)

type Template struct {
	attributes map[string]interface{}
	source     string
}

func (t *Template) render() string {
	// ERB.new(t.source).result(binding)
	return t.source
}

type Entries struct {
	Attributes map[string]interface{} `json:"attributes"`
	Servers    map[string][]string    `json:"servers"`
}

func layout() map[string]Entries {
	input, err := ioutil.ReadFile(path.Join(home, "layout.json"))

	if err != nil {
		fmt.Println(err)
	}

	var e map[string]Entries

	json.Unmarshal(input, &e)

	return e
}

func run(server, recipe, command string, attributes map[string]interface{}) (string, int) {
	template_path := path.Join(home, "recipe", recipe, command)

	if _, err := os.Stat(template_path); os.IsNotExist(err) {
		fmt.Print("\033[01;33mMISSING\033[00m")
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

	var stderrb bytes.Buffer
	var stdoutb bytes.Buffer

	cmd.Stderr = &stderrb
	cmd.Stdout = &stdoutb

	err := cmd.Run()

	if err != nil {
		fmt.Print(err)
		return stdoutb.String() + stderrb.String(), 1
	}

	// out, status = Open3.capture2e("ssh -T -F #{path("ssh_config")} #{server}", :stdin_data => script)

	return stdoutb.String() + stderrb.String(), 0
}

var home string

func init() {
	home, err := os.Getwd()

	if err != nil {
		fmt.Println("unable to get pwd", home, err)
		os.Exit(1)
	}
}

func main() {
	// parse commands
	// check layout
	//

	var command = flag.String("c", "", "command")
	var directory = flag.String("d", ".", "directory")
	var environment = flag.String("e", "", "environment")
	var quiet = flag.Bool("q", false, "quiet mode")
	var server = flag.String("s", "", "server")
	var verbose = flag.Bool("v", false, "verbose mode")

	flag.Parse()

	fmt.Println(home, *directory, *quiet, *verbose)

	e := layout()

	layout := e[*environment]

	var servers []string

	for k, _ := range layout.Servers {
		servers = append(servers, k)
	}

	exit_status := 0

	for _, v := range servers {
		recipes := layout.Servers[v]

		fmt.Print(v)

		for _, recipe := range recipes {
			fmt.Printf("  %s: ", recipe)

			filename := path.Join(home, "recipe", recipe)

			if _, err := os.Stat(filename); os.IsNotExist(err) {
				fmt.Printf("unable to locate: %s\n", filename)
				os.Exit(1)
			}

			stdout, status := run(v, recipe, *command, layout.Attributes)

			switch status {
			case -1: // nil is better? negative error codes?
				fmt.Print("?")
			case 0:
				fmt.Print("\033[01;32mOK\033[00m")
			default:
				fmt.Print("\033[01;31mERROR\033[00m")
				exit_status = 1
				break
			}

			if len(stdout) > 0 {
				fmt.Fprintf(os.Stderr, " %s\n", stdout)
			}
		}
	}

	fmt.Println("\033[01;32mDONE\033[00m")
	os.Exit(exit_status)
}
