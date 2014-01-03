package layouts

import (
	"encoding/json"
	"io/ioutil"
	"path"
)

type Environment struct {
	Attributes map[string]interface{} `json:"attributes"`
	Servers    map[string][]string    `json:"servers"`
}

type Layout map[string]Environment

func Read(directory string) ([]byte, error) {
	input, err := ioutil.ReadFile(path.Join(directory, "layout.json"))

	return input, err
}

func (layout *Layout) Parse(input []byte) error {
	err := json.Unmarshal(input, &layout)

	return err
}
