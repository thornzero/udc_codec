package udc

import (
	"strings"
)

func search(flat map[string]*Node, term string) []*Node {
	term = strings.ToLower(term)
	var results []*Node
	for _, node := range flat {
		if strings.Contains(strings.ToLower(node.Title), term) {
			results = append(results, node)
		}
	}
	return results
}
