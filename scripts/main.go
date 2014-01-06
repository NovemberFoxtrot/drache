package scripts

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"text/template"
)

type Script struct {
	Attributes  map[string]interface{}
	Command     string
	Directory   string
	Environment string
	Output      string
	Server      string
	Name        string
	Status      int
}

func (script *Script) location() string {
	return path.Join(script.Directory, "scripts", script.Name, script.Command)
}

func (script *Script) missing() bool {
	if _, err := os.Stat(script.location()); os.IsNotExist(err) {
		return true
	}

	return false
}

func (script *Script) source() string {
	if source, err := ioutil.ReadFile(script.location()); err != nil {
		panic(err)
	} else {
		tmpl, err := template.New(script.location()).Parse(string(source))

		if err != nil {
			fmt.Println("parsing: %s", err)
		}

		b := new(bytes.Buffer)

		err = tmpl.Execute(b, script.Attributes)

		if err != nil {
			fmt.Println("execution: %s", err)
		}

		return string(b.String())
	}
}

func (script *Script) Run() {
	if script.missing() {
		script.Output = "\033[01;33mMISSING\033[00m unable to locate: " + script.location()
		script.Status = 1
		return
	}

	out, status := script.ssh(script.Server, script.source())
	script.Output = out
	script.Status = status
}

func (script *Script) sshConfigLocation() string {
	return path.Join(script.Directory, "ssh_config")
}

func (script *Script) ssh(server, source string) (string, int) {
	cmd := exec.Command("ssh", "-T", "-F", script.sshConfigLocation(), server, source)

	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Print(err)
		return string(output), 1
	}

	return string(output), 0
}
