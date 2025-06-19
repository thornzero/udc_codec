package udc

import (
	"fmt"
	"regexp"
	"strings"
)

// Split composite codes like 621.3:681.5(075)
func parseComposite(code string) ([]string, error) {
	code = strings.ReplaceAll(code, " ", "")
	re := regexp.MustCompile(`(\d+(\.\d+)*|\(\d+\)|\([^\)]+\))`)
	parts := re.FindAllString(code, -1)
	if len(parts) == 0 {
		return nil, fmt.Errorf("invalid UDC code: %s", code)
	}
	return parts, nil
}
