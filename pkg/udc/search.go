package udc

import (
	"strings"
)

func search(flat map[string]*TreeNode, term string) []*TreeNode {
	term = strings.ToLower(term)
	var results []*TreeNode
	for _, node := range flat {
		if strings.Contains(strings.ToLower(node.Title), term) {
			results = append(results, node)
		}
	}
	return results
}
