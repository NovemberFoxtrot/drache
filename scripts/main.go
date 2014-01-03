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

func (b *Script) exec(server, recipe string) (string, int) {
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

func (b *Script) Run() {
	/*
		stdout, status := b.exec(server, script)

		if status != 0 {
			fmt.Fprintf(os.Stderr, "\033[01;31mERROR\033[00m\n %s\n", stdout)
			b.Status = 1
			break

		}

		fmt.Print("\033[01;32mOK\033[00m\n")
	*/
}
