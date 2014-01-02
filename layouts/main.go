package layouts

import (
	"encoding/json"
	"io/ioutil"
	"path"
)

type Layout struct {
	Attributes map[string]interface{} `json:"attributes"`
	Servers    map[string][]string    `json:"servers"`
}

func Read() ([]byte, error) {
	input, err := ioutil.ReadFile(path.Join(".", "layout.json"))

	return input, err
}

func (layout *Layout) Parse(input []byte) error {
	err := json.Unmarshal(input, &layout)

	return err
}
