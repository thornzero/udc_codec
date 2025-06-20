package pipeline

import (
	"fmt"
)

func GenerateFullTag(entry BOMEntry) string {
	tag := fmt.Sprintf("%s-%s%s", entry.SystemCode, entry.FunctionCode, entry.EquipmentID)
	return tag
}
