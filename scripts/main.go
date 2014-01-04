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
	Environment string
	Status      int
}

func (b *Script) Run(server, recipe string) (string, int) {
	scriptPath := path.Join(".", "recipe", recipe, b.Command)

	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return "\033[01;33mMISSING\033[00m unable to locate: " + scriptPath, 0
	}

	if source, err := ioutil.ReadFile(scriptPath); err != nil {
		return "unable to read file: " + scriptPath, 1
	} else {
		out, status := ssh(server, string(source))
		return out, status
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
