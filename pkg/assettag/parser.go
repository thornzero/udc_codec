package assettag

import (
	"fmt"
	"regexp"
)

// Example tag: POL-LT1001-A

var tagRegex = regexp.MustCompile(`^([A-Z]{3})-([A-Z]{1,5})(\d{3,5})(-[A-Z])?$`)

func ParseTag(tag string) (*Tag, error) {
	parts := tagRegex.FindStringSubmatch(tag)
	if len(parts) == 0 {
		return nil, fmt.Errorf("invalid tag format")
	}
	return &Tag{
		SystemCode:   parts[1],
		FunctionCode: parts[2],
		EquipmentID:  parts[3],
		InstrumentID: parts[4],
	}, nil
}
