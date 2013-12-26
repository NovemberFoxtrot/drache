package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	_ "io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
)

type Template struct {
	source     string
	attributes map[string]interface{}
}

func (t *Template) initialize(source string) {
	t.source = source
}

func (t *Template) render() string {
	// ERB.new(t.source).result(binding)
	return t.source
}

func layout() {
	// @layout ||= JSON.parse(File.read(path("layout.json")))
}

func run(server, recipe, command string, attributes map[string]interface{}) (string, int) {
	template_path := path.Join(home, "recipe", recipe, command)

	if _, err := os.Stat(template_path); os.IsNotExist(err) {
		return "unable to locate: " + template_path, 1
	}

	source, err := ioutil.ReadFile(template_path)

	if err != nil {
		return "unable to read file: " + template_path, 1
	}

	template := new(Template)
	template.source = string(source)
	template.attributes = attributes
	ssh(server, template.render())

	return "", 0
}

func local() {
	// runs recipe locally
}

func telnet() {
	// uses telnet
}

func ssh(server, script string) (string, int) {
	cmd := exec.Command(script[:len(script)-1]) // TODO

	var stderrb bytes.Buffer
	var stdoutb bytes.Buffer

	cmd.Stdout = &stdoutb
	cmd.Stderr = &stderrb
	err := cmd.Run()

	if err != nil {
		fmt.Println(err)
		return "Errrrrr", 1
	}

	fmt.Println("stdout", stdoutb.String(), "stderr", stderrb.String(), "err", err)

	// out, status = Open3.capture2e("ssh -T -F #{path("ssh_config")} #{server}", :stdin_data => script)
	// return out, status.exitstatus

	return "", 0
}

type out struct {
}

func (o *out) server( /* name */) {
	fmt.Print("" /* name */)
}

func (o *out) error() {
	fmt.Print("\033[01;31mERROR\033[00m")
}

func (o *out) ok() {
	fmt.Print("\033[01;32mOK\033[00m")
}

func (o *out) missing() {
	fmt.Print("\033[01;33mMISSING\033[00m")
}

func (o *out) done() {
	fmt.Print("\033[01;32mDONE\033[00m")
}

func (o *out) unknown() {
	fmt.Print("?")
}

var home string

func init() {
	home, err := os.Getwd()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(home)
}

func main() {
	var directory = flag.String("d", ".", "directory") // recipe / layout.json root
	var quiet = flag.Bool("q", false, "quiet mode")
	var verbose = flag.Bool("v", false, "verbose mode")
	var environment = flag.String("e", "development", "environment")

	flag.Parse()

	fmt.Println(home, *directory, *quiet, *verbose, *environment)

	input, err := ioutil.ReadFile("layout.json")

	if err != nil {
		fmt.Println(err)
	}

	type Entries struct {
		Servers    map[string][]string    `json:"servers"`
		Attributes map[string]interface{} `json:"attributes"`
	}

	var e map[string]Entries

	json.Unmarshal(input, &e)

	switch os.Args[1] {
	case "run":
		command := os.Args[2]

		environment := strings.Split(os.Args[3], ":")[0]

		server := ""

		if len(strings.Split(os.Args[3], ":")) > 1 {
			server = strings.Split(os.Args[3], ":")[1]
		}

		layout := e[environment]

		var servers []string

		if len(strings.Split(server, ",")) > 1 {
			for _, v := range strings.Split(server, ",") {
				servers = append(servers, v)
			}
		} else {
			for k, _ := range layout.Servers {
				servers = append(servers, k)
			}
		}

		attributes := layout.Attributes

		exit_status := 0

		fmt.Println(command, environment, server, layout, servers, attributes, exit_status)

		for _, v := range servers {
			recipes := layout.Servers[v]

			o := new(out)

			o.server()
			fmt.Println(v)

			for _, recipe := range recipes {
				fmt.Printf("  %s: ", recipe)

				filename := path.Join(home, "recipe", recipe)

				if _, err := os.Stat(filename); os.IsNotExist(err) {
					fmt.Printf("unable to locate: %s\n", filename)
					os.Exit(1)
				}

				stdout, status := run(server, recipe, command, attributes)

				switch status {
				case -1: // nil
					o.unknown()
					if len(stdout) > 0 {
						fmt.Fprintf(os.Stderr, " %s\n", stdout)
					}
				case 0:
					o.ok()
					if len(stdout) > 0 {
						fmt.Fprintf(os.Stderr, " %s\n", stdout)
					}
				default:
					o.error()
					if len(stdout) > 0 {
						fmt.Fprintf(os.Stderr, " %s\n", stdout)
					}
					exit_status = 1
					break
				}
			}
		}
		os.Exit(exit_status)
	}
}
