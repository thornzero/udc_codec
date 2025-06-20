package aggregator

import (
	"os"

	"gopkg.in/yaml.v3"
)

func LoadAggregatedDatabase(filename string) (*AggregatedDatabase, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var db AggregatedDatabase
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&db); err != nil {
		return nil, err
	}
	return &db, nil
}
