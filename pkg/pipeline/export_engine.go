package pipeline

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func LoadExportedTags(filename string) ([]ExportRecord, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var records []ExportRecord
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&records); err != nil {
		return nil, err
	}
	return records, nil
}

func ExportMarkdown(project string, entries []ExportRecord) error {
	f, err := os.Create(fmt.Sprintf("data/%s_doc.md", project))
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintf(f, "# Project: %s\n\n", project)
	fmt.Fprintln(f, "| FullTag | System | Description | UDC |")
	fmt.Fprintln(f, "|---------|--------|-------------|-----|")
	for _, rec := range entries {
		fmt.Fprintf(f, "| %s | %s | %s | %s |\n", rec.FullTag, rec.SystemName, rec.Description, rec.UDCCode)
	}
	return nil
}
