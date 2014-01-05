package scripts

import (
	"testing"
)

func TestLocation(t *testing.T) {
	script := &Script{Command: "install", Name: "go", Directory: "."}

	actual := script.location()
	if actual != "scripts/go/install" {
		t.Errorf("%s", actual)
	}
}

func TestMissing(t *testing.T) {
	script := &Script{Command: "install", Name: "go", Directory: "."}

	actual := script.missing()
	if actual != true {
		t.Errorf("%s", actual)
	}
}
