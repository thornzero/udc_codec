package main

import (
	"fmt"
	"log"

	"github.com/thornzero/udc_codec/pkg/aggregator"
	"github.com/thornzero/udc_codec/pkg/pipeline"
	"github.com/thornzero/udc_codec/pkg/udc"
)

func main() {
	// Load UDC codec
	udcCodec, err := udc.LoadCodec("data/udc_full.yaml")
	if err != nil {
		log.Fatalf("UDC load failed: %v", err)
	}

	// Load Aggregator
	agg, err := aggregator.LoadAggregatedDatabase("data/aggregated_master.yaml")
	if err != nil {
		log.Fatalf("Aggregator load failed: %v", err)
	}

	// Load Project BOM
	bom, err := pipeline.LoadBOM("data/project_bom.yaml")
	if err != nil {
		log.Fatalf("BOM load failed: %v", err)
	}

	// Prepare validator
	validator := &pipeline.Validator{
		Aggregator: agg,
		UDC:        udcCodec,
	}

	// Full pipeline process
	var exportRecords []pipeline.ExportRecord
	for _, entry := range bom.Entries {
		// Validate entry
		if err := validator.ValidateEntry(entry); err != nil {
			log.Fatalf("Validation failed for entry %+v: %v", entry, err)
		}

		tag := pipeline.GenerateFullTag(entry)
		system := agg.LookupSystem(entry.SystemCode)

		exportRecords = append(exportRecords, pipeline.ExportRecord{
			FullTag:     tag,
			SystemName:  system.SystemName,
			Description: entry.Description,
			UDCCode:     entry.UDCCode,
		})
	}

	outputFile := fmt.Sprintf("data/%s_taglist.yaml", bom.ProjectName)
	if err := pipeline.ExportTagList(exportRecords, outputFile); err != nil {
		log.Fatalf("Export failed: %v", err)
	}

	fmt.Println("✅ Pipeline complete!")
	fmt.Printf("Exported tag list to: %s\n", outputFile)
}
