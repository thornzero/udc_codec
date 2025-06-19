package udc

import (
	"fmt"
)

func validateComposite(code string, flat map[string]*TreeNode) error {
	parts, err := parseComposite(code)
	if err != nil {
		return err
	}
	for _, p := range parts {
		if _, ok := flat[p]; !ok {
			return fmt.Errorf("invalid code part: %s", p)
		}
	}
	return nil
}
