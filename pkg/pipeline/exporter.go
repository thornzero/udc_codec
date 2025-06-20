package pipeline

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ExportRecord struct {
	FullTag     string `yaml:"full_tag"`
	SystemName  string `yaml:"system_name"`
	Description string `yaml:"description"`
	UDCCode     string `yaml:"udc_code,omitempty"`
}

func ExportTagList(entries []ExportRecord, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := yaml.NewEncoder(f)
	encoder.SetIndent(2)
	return encoder.Encode(entries)
}
