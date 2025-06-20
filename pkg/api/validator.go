package api

import (
	"fmt"

	"github.com/thornzero/udc_codec/pkg/pipeline"
)

// ValidateAPITag checks for required fields and delegates deep validation to pipeline.Validator.
func ValidateAPITag(tag APITag, validator *pipeline.Validator) error {
	if tag.SystemCode == "" {
		return fmt.Errorf("system_code is required")
	}
	if tag.FunctionCode == "" {
		return fmt.Errorf("function_code is required")
	}
	if tag.EquipmentID == "" {
		return fmt.Errorf("equipment_id is required")
	}
	// Optionally check for description or UDCCode presence/format here

	// Convert APITag to pipeline.BOMEntry for deep validation
	entry := pipeline.BOMEntry{
		SystemCode:   tag.SystemCode,
		FunctionCode: tag.FunctionCode,
		EquipmentID:  tag.EquipmentID,
		UDCCode:      tag.UDCCode,
		Description:  tag.Description,
	}
	return validator.ValidateEntry(entry)
}
