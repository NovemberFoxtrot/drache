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

func ReadLayout() ([]byte, error) {
	input, err := ioutil.ReadFile(path.Join(".", "layout.json"))

	return input, err
}

func (layout *Layout) ParseLayout(input []byte) error {
	err := json.Unmarshal(input, &layout)

	return err
}
