package main

import (
	"fmt"
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

func run( /* server, recipe, command, attributes */ /* ={} */) {
	// template_path = path("recipes", recipe, "#{command}.erb")

	// if File.exists?(template_path)
	//  source = File.read(template_path)
	//  template = Template.new(source)
	//  template.attributes = attributes
	//  ssh(server, template.render)
	// end
}

func telnet() {
}

func ssh(server, script string) {
	// out, status = Open3.capture2e("ssh -T -F #{path("ssh_config")} #{server}", :stdin_data => script)
	// [out, status.exitstatus]
}

/*
commands = Clap.run ARGV,
  "-q" => lambda {
    $verbosity = 0
  },
  "-v" => lambda {
    $verbosity = 2
  },
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
	// $home = Dir.pwd

	switch os.Args[1] {
	case "run":
		command := os.Args[2]

		environment := strings.Split(os.Args[3], ":")[0]
		servers := strings.Split(os.Args[3], ":")[1]

		fmt.Println(command, environment, servers)
	}

	/*
	   case ARGV.shift
	   when "run" then
	     command = ARGV.shift
	     environment, servers = ARGV.shift.split(":")

	     environment = layout[environment]

	     servers = if servers
	                 servers.split(",")
	               else
	                 environment["servers"].keys
	               end

	     attributes = environment["attributes"] || {}

	     exit_status = 0

	     servers.each do |server|
	       recipes = environment["servers"][server]

	       out.server(server)

	       recipes.each do |recipe|
	         print "  #{recipe}: "

	         if File.exists?(path("recipes", recipe))
	           stdout, status = run(server, recipe, command, attributes)

	           case status
	           when nil
	             out.unknown
	           when 0
	             out.ok
	             $stderr.print stdout if $verbosity >= 2
	           else
	             out.error
	             $stderr.print stdout if $verbosity >= 1
	             exit_status = 1
	             break
	           end
	         else
	           out.unknown
	           exit 1
	         end
	       end
	     end

	     exit exit_status
	   end
	*/
}
