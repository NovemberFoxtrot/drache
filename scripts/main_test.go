package scripts

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	book := &Book{command: "install", environment: "testing", status: 0}

	if book.command != "install" {
		t.Errorf("Dude")
	}
}

func TestNoLayout(t *testing.T) {
	/*
		currentDir, err := os.Getwd()

		if err != nil {
			panic(err)
		}
	*/

	dir, err := ioutil.TempDir(".", "test_layout")

	if err != nil {
		panic(err)
	}

	defer func() {
		err := os.RemoveAll(dir)

		if err != nil {
			panic(err)
		}

		/*
			err = os.Chdir(dir)

			if err != nil {
				panic(err)
			}
		*/
	}()

	/*
		err = os.Chdir(dir)

		if err != nil {
			panic(err)
		}
	*/

	cmd := exec.Command("./drache", "install", "development")
	output, err := cmd.CombinedOutput()

	if err != nil {
		panic(err)
	}

	if strings.Contains(string(output), "OK") != true {
		t.Errorf("Layout issuei: %s", string(output))
	}

	// assert exit status
}
