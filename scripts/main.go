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

func (script *Script) Run() {
	scriptPath := path.Join(script.Directory, "recipe", script.Name, script.Command)

	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		script.Output = "\033[01;33mMISSING\033[00m unable to locate: " + scriptPath
		script.Status = 1
	}

	if source, err := ioutil.ReadFile(scriptPath); err != nil {
		script.Output = "unable to read file: " + scriptPath
		script.Status = 1
	} else {
		out, status := ssh(script.Server, string(source))
		script.Output = out
		script.Status = status
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
