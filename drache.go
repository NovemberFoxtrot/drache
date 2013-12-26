package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Template struct {
	// attr_accessor :attributes
	source string
}

func (t *Template) initialize(source string) {
	t.source = source
}

func (t *Template) render() {
	// ERB.new(t.source).result(binding)
}

func path( /* *parts */) {
	// File.join($home, *parts)
}

func layout() {
	// @layout ||= JSON.parse(File.read(path("layout.json")))
}

func run(server, recipe, command string, attributes map[string]interface{}) (string, int) {
	// template_path = path("recipes", recipe, "#{command}.erb")

	// if File.exists?(template_path)
	//  source = File.read(template_path)
	//  template = Template.new(source)
	//  template.attributes = attributes
	//  ssh(server, template.render)
	// end

	return "temp", 0
}

func local() {
	// runs recipe locally
}

func telnet() {
	// uses telnet
}

func ssh(server, script string) {
	// out, status = Open3.capture2e("ssh -T -F #{path("ssh_config")} #{server}", :stdin_data => script)
	// [out, status.exitstatus]
}

/*
commands = Clap.run ARGV,
  "-d" => lambda {|path|
    $home = File.join(Dir.pwd, path)
  }

unless File.exists?(path("layout.json"))
  $stderr.puts "Couldn't find `layout.json`"
  exit 1
end
*/

type out struct {
}

func (o *out) server( /* name */) {
	fmt.Println("" /* name */)
}

func (o *out) error() {
	fmt.Println("\033[01;31mERROR\033[00m")
}

func (o *out) ok() {
	fmt.Println("\033[01;32mOK\033[00m")
}

func (o *out) missing() {
	fmt.Println("\033[01;33mMISSING\033[00m")
}

func (o *out) done() {
	fmt.Println("\033[01;32mDONE\033[00m")
}

func (o *out) unknown() {
	fmt.Println("?")
}

func main() {
	home, err := os.Getwd()

	var path = flag.String("d", ".", "path")
	var quiet = flag.Bool("q", false, "quiet mode")
	var verbose = flag.Bool("v", false, "verbose mode")
	var environment = flag.String("e", "development", "environment")

	flag.Parse()

	fmt.Println(home, *path, *quiet, *verbose, *environment)

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

				filename := home + "/recipe/" + recipe

				if _, err := os.Stat(filename); os.IsNotExist(err) {
					fmt.Printf("unable to locate: %s\n", filename)
					os.Exit(1)
				}

				stdout, status := run(server, recipe, command, attributes)

				switch status {
				case -1: // nil
					o.unknown()
				case 0:
					o.ok()
					if *verbose == true {
						fmt.Fprintf(os.Stderr, "%s\n", stdout)
					}
				default:
					o.error()
					if *verbose == true {
						fmt.Fprintf(os.Stderr, "%s\n", stdout)
					}
					exit_status = 1
					break
				}
			}
		}
		os.Exit(exit_status)
	}
}
