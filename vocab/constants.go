package vocab

import (
	"encoding/json"
	"strings"
)

type Constant struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Kind  string `json:"kind"`
}

func Constants(prefix string) ([]Constant, error) {
	data, err := Bundled("constants.json")
	if err != nil {
		return nil, err
	}
	var all []Constant
	if err := json.Unmarshal(data, &all); err != nil {
		return nil, err
	}
	if prefix == "" {
		return all, nil
	}
	out := make([]Constant, 0)
	for _, c := range all {
		if strings.HasPrefix(c.Name, prefix) {
			out = append(out, c)
		}
	}
	return out, nil
}
