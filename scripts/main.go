package scripts

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
)

type Script struct {
	Command     string
	Directory   string
	Environment string
	Output      string
	Server      string
	Name        string
	Status      int
}

func (script *Script) location() string {
	return path.Join(script.Directory, "recipe", script.Name, script.Command)
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
		return string(source)
	}
}

func (script *Script) Run() {
	if script.missing() {
		script.Output = "\033[01;33mMISSING\033[00m unable to locate: " + script.location()
		script.Status = 1
		return
	}

	out, status := ssh(script.Server, script.source())
	script.Output = out
	script.Status = status
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
