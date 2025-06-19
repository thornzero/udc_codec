package main

import (
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/thornzero/udc_codec/pkg/assettag"
	"github.com/thornzero/udc_codec/pkg/config"
	"github.com/thornzero/udc_codec/pkg/db"
	"github.com/thornzero/udc_codec/pkg/udc"
)

func openFile(filename string) (*os.File, error) {
	cfg := config.Load()
	return os.Open(cfg.Path(filename))
}

// YAML structure loaders

func loadISA(filename string) (map[string]string, error) {
	f, err := openFile(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var data map[string]string
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

func loadIEC81346(filename string) (map[string]string, error) {
	f, err := openFile(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var data map[string]string
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

func main() {
	// Load YAML data
	codec, err := udc.LoadCodec("data/udc_full.yaml")
	if err != nil {
		panic(err)
	}

	isa, err := loadISA("data/isa_prefix.yaml")
	if err != nil {
		panic(err)
	}

	systems, err := loadIEC81346("data/81346_systems.yaml")
	if err != nil {
		panic(err)
	}

	resolver := &assettag.Resolver{
		UDC:      codec,
		ISA:      isa,
		IEC81346: systems,
	}

	// Open DB
	store, err := db.OpenDB("tags.db")
	if err != nil {
		panic(err)
	}
	if err := store.Migrate(); err != nil {
		panic(err)
	}

	// Load initial asset tags (you supply this file)
	f, err := openFile("data/asset_tags.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	var tags []string
	if err := decoder.Decode(&tags); err != nil {
		panic(err)
	}

	for _, fulltag := range tags {
		tag, err := assettag.ParseTag(fulltag)
		if err != nil {
			fmt.Printf("Skipping invalid tag: %s (%v)\n", fulltag, err)
			continue
		}

		if err := resolver.ValidateTag(tag); err != nil {
			fmt.Printf("Validation failed: %s (%v)\n", fulltag, err)
			continue
		}

		desc := resolver.DescribeTag(tag)
		rec := db.TagRecord{
			FullTag:      fulltag,
			SystemCode:   tag.SystemCode,
			EquipmentID:  tag.EquipmentID,
			InstrumentID: tag.InstrumentID,
			FunctionCode: tag.FunctionCode,
			UDCCode:      tag.UDCCode,
			Description:  desc,
		}
		if err := store.InsertTag(&rec); err != nil {
			fmt.Printf("DB insert error for %s: %v\n", fulltag, err)
		}
	}

	fmt.Println("✅ Bootstrap complete.")
}
