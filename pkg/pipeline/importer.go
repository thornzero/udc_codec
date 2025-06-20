package pipeline

import (
	"os"

	"gopkg.in/yaml.v3"
)

func LoadBOM(filename string) (*ProjectBOM, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var bom ProjectBOM
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&bom); err != nil {
		return nil, err
	}
	return &bom, nil
}
