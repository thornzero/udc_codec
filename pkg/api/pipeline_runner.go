package api

import (
	"fmt"

	"github.com/thornzero/udc_codec/pkg/aggregator"
	"github.com/thornzero/udc_codec/pkg/db"
	"github.com/thornzero/udc_codec/pkg/pipeline"
	"github.com/thornzero/udc_codec/pkg/udc"
)

func runFullPipeline(projectName, bomFile string) error {
	udcCodec, err := udc.LoadCodec("data/udc_full.yaml")
	if err != nil {
		return err
	}
	agg, err := aggregator.LoadAggregatedDatabase("data/aggregated_master.yaml")
	if err != nil {
		return err
	}
	bom, err := pipeline.LoadBOM(bomFile)
	if err != nil {
		return err
	}

	validator := &pipeline.Validator{
		Aggregator: agg,
		UDC:        udcCodec,
	}

	var exportRecords []pipeline.ExportRecord
	for _, entry := range bom.Entries {
		if err := validator.ValidateEntry(entry); err != nil {
			return fmt.Errorf("validation failed: %v", err)
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

	outputFile := fmt.Sprintf("data/%s_taglist.yaml", projectName)
	if err := pipeline.ExportTagList(exportRecords, outputFile); err != nil {
		return err
	}

	// Insert into DB registry
	store, err := db.OpenDB("tags.db")
	if err != nil {
		return err
	}
	if err := store.Migrate(); err != nil {
		return err
	}
	_, err = store.InsertProject(db.ProjectRecord{
		ProjectName: projectName,
		FullBOMFile: bomFile,
		Validated:   true,
	})
	return err
}
