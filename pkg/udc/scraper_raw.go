package udc

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/chromedp/chromedp"
)

// Debug flag to control logging verbosity
var DebugMode = false

// SetDebugMode enables or disables debug logging
func SetDebugMode(enabled bool) {
	DebugMode = enabled
}

// IsDebugMode returns the current debug mode status
func IsDebugMode() bool {
	return DebugMode
}

type RawNode struct {
	ID       string
	Parent   string
	Code     string
	Title    string
	Children []*RawNode
}

func ScrapeRawTree(ctx context.Context) ([]*RawNode, error) {
	var rawHTML string
	err := chromedp.Run(ctx,
		chromedp.OuterHTML("#classtree", &rawHTML),
	)
	if err != nil {
		return nil, err
	}

	if DebugMode {
		fmt.Printf("[DEBUG] Length of rawHTML: %d\n", len(rawHTML))
	}

	parsed := parseRawHTML(rawHTML)
	if DebugMode {
		fmt.Printf("[DEBUG] Number of nodes parsed: %d\n", len(parsed))
		fmt.Printf("[DEBUG] Sample parsed nodes:\n")
		for i, node := range parsed {
			if i >= 5 { // Show first 5 nodes as sample
				break
			}
			fmt.Printf("[DEBUG]   %s: %s - %s (parent: %s)\n", node.ID, node.Code, node.Title, node.Parent)
		}
	}

	root := buildRawHierarchy(parsed)
	if DebugMode {
		fmt.Printf("[DEBUG] Number of root nodes: %d\n", len(root))
		fmt.Printf("[DEBUG] Root nodes:\n")
		for _, r := range root {
			fmt.Printf("[DEBUG]   %s: %s - %s (%d children)\n", r.ID, r.Code, r.Title, len(r.Children))
		}
	}
	return root, nil
}

func parseRawHTML(html string) []*RawNode {
	// Fixed regex that matches the test HTML format exactly
	re := regexp.MustCompile(`d\.add\(\s*(\d+),\s*(-\d|\d+),\s*\'(.*?)\'\,\s*\'[^\']*?&nbsp;&nbsp;(.*?)\'`)
	matches := re.FindAllStringSubmatch(html, -1)

	var nodes []*RawNode
	for _, m := range matches {
		code := strings.TrimSpace(m[3])
		if code == "-" || code == "--" || code == "---" || code == "----" {
			continue
		}
		node := &RawNode{
			ID:     m[1],
			Parent: m[2],
			Code:   code,
			Title:  strings.TrimSpace(m[4]),
		}
		nodes = append(nodes, node)
	}
	return nodes
}

func buildRawHierarchy(nodes []*RawNode) []*RawNode {
	// Create maps for efficient lookup
	codeMap := make(map[string]*RawNode)
	idMap := make(map[string]*RawNode)

	// First pass: build maps
	for _, n := range nodes {
		codeMap[n.Code] = n
		idMap[n.ID] = n
	}

	if DebugMode {
		fmt.Printf("[DEBUG] Building hierarchy for %d nodes\n", len(nodes))
	}

	// Second pass: build hierarchy based on UDC code structure
	for _, n := range nodes {
		parentCode := findParentCode(n.Code)
		if parentCode != "" {
			if parent, ok := codeMap[parentCode]; ok {
				parent.Children = append(parent.Children, n)
				if DebugMode {
					fmt.Printf("[DEBUG] %s → %s (code-based parent)\n", n.Code, parentCode)
				}
				continue
			} else if DebugMode {
				fmt.Printf("[DEBUG] %s → %s (parent not found in codeMap)\n", n.Code, parentCode)
			}
		}

		// Fallback to website's parent ID if code-based parent not found
		// But only for nodes that shouldn't be roots
		if n.Parent != "0" && n.Parent != "-1" && !shouldBeRoot(n.Code) {
			if parent, ok := idMap[n.Parent]; ok {
				parent.Children = append(parent.Children, n)
				if DebugMode {
					fmt.Printf("[DEBUG] %s → %s (fallback to website parent ID)\n", n.Code, parent.Code)
				}
				continue
			} else if DebugMode {
				fmt.Printf("[DEBUG] %s → parent ID %s (parent not found in idMap)\n", n.Code, n.Parent)
			}
		}

		// If no parent found, it's a root node
		if DebugMode {
			fmt.Printf("[DEBUG] %s → ROOT (no parent found)\n", n.Code)
		}
	}

	// Collect root nodes
	var roots []*RawNode
	for _, n := range nodes {
		// Check if this node has no parent assigned
		hasParent := false
		for _, other := range nodes {
			for _, child := range other.Children {
				if child == n {
					hasParent = true
					break
				}
			}
			if hasParent {
				break
			}
		}

		if !hasParent {
			roots = append(roots, n)
		}
	}

	if DebugMode {
		fmt.Printf("[DEBUG] Final root count: %d\n", len(roots))
	}

	return roots
}

// findParentCode determines the parent UDC code based on the current code
func findParentCode(code string) string {
	if code == "" {
		return ""
	}

	// Handle auxiliary tables
	if strings.HasPrefix(code, "=") {
		return findAuxiliaryParent(code)
	}
	if strings.HasPrefix(code, "(") {
		return findAuxiliaryParent(code)
	}
	if strings.HasPrefix(code, "-") {
		return findAuxiliaryParent(code)
	}
	if strings.HasPrefix(code, "+") || strings.HasPrefix(code, "/") ||
		strings.HasPrefix(code, ":") || strings.HasPrefix(code, "[]") ||
		strings.HasPrefix(code, "*") || strings.HasPrefix(code, "A/Z") {
		return "TOP" // These are top-level auxiliary signs
	}

	// Handle main tables (0-9)
	return findMainTableParent(code)
}

// findAuxiliaryParent finds the parent for auxiliary table codes
func findAuxiliaryParent(code string) string {
	if code == "" {
		return ""
	}

	// Handle common auxiliaries of language (=...)
	if strings.HasPrefix(code, "=...") {
		if strings.Contains(code, "`") {
			// Handle special auxiliary subdivisions like =...`01/`08
			parts := strings.Split(code, "`")
			if len(parts) > 1 {
				return "=..."
			}
		}
		return "TOP" // =... is a top-level auxiliary
	}

	// Handle language codes (=1, =11, =111, etc.)
	if strings.HasPrefix(code, "=") && len(code) > 1 {
		// Remove the = prefix
		langCode := code[1:]

		// Handle special cases first
		if langCode == "00" || langCode == "030" {
			return "=..."
		}

		// Handle single-digit language codes (top-level)
		if len(langCode) == 1 && isNumeric(langCode) {
			return "TOP"
		}

		// Handle numeric language codes
		if isNumeric(langCode) {
			return findNumericParent(langCode, "=", "")
		}
	}

	// Handle place auxiliaries ((1), (1-44), etc.)
	if strings.HasPrefix(code, "(") && strings.HasSuffix(code, ")") {
		placeCode := code[1 : len(code)-1]

		// Handle special cases like (=01)
		if strings.HasPrefix(placeCode, "=") {
			return "(=...)"
		}

		// Handle single-digit place codes (top-level)
		if len(placeCode) == 1 && isNumeric(placeCode) {
			return "TOP"
		}

		// Handle numeric place codes - for place codes, we need to handle them differently
		// (540) should be a child of (5), not (54)
		if isNumeric(placeCode) {
			if len(placeCode) > 1 {
				// For place codes, take the first digit as parent
				firstDigit := string(placeCode[0])
				return "(" + firstDigit + ")"
			}
		}
	}

	// Handle form auxiliaries (-0, -058.6, etc.)
	if strings.HasPrefix(code, "-") {
		formCode := code[1:]

		// Handle special cases
		if formCode == "0" {
			return "TOP" // -0 is a top-level auxiliary
		}

		// Handle numeric form codes
		if isNumeric(formCode) {
			return findNumericParent(formCode, "-", "")
		}
	}

	return ""
}

// findMainTableParent finds the parent for main table codes (0-9)
func findMainTableParent(code string) string {
	if code == "" {
		return ""
	}

	// Handle special cases
	if code == "0" || code == "1" || code == "2" || code == "3" ||
		code == "4" || code == "5" || code == "6" || code == "7" ||
		code == "8" || code == "9" {
		return "TOP" // Main table divisions are children of TOP
	}

	// Handle numeric codes with dots (001, 001.1, etc.)
	if strings.Contains(code, ".") {
		parts := strings.Split(code, ".")
		if len(parts) > 1 {
			// Remove the last part to get parent
			parentParts := parts[:len(parts)-1]
			return strings.Join(parentParts, ".")
		}
	}

	// Handle numeric codes without dots (00, 000, 01, etc.)
	if isNumeric(code) {
		// Remove last digit to get parent
		if len(code) > 1 {
			return code[:len(code)-1]
		}
	}

	return ""
}

// findNumericParent finds the parent for numeric codes with a prefix and optional suffix
func findNumericParent(code, prefix, suffix string) string {
	if code == "" {
		return ""
	}

	// Handle special ranges like 01/08
	if strings.Contains(code, "/") {
		parts := strings.Split(code, "/")
		if len(parts) > 1 {
			// For ranges, the parent is the prefix part
			return prefix + parts[0] + suffix
		}
	}

	// Handle numeric codes
	if isNumeric(code) {
		if len(code) > 1 {
			// Remove last digit to get parent, but preserve dots
			lastChar := code[len(code)-1]
			if lastChar >= '0' && lastChar <= '9' {
				parentCode := code[:len(code)-1]
				// Remove trailing dot if present
				parentCode = strings.TrimSuffix(parentCode, ".")
				return prefix + parentCode + suffix
			}
		}
	}

	return ""
}

// isNumeric checks if a string contains only digits and dots
func isNumeric(s string) bool {
	for _, r := range s {
		if r != '.' && (r < '0' || r > '9') {
			return false
		}
	}
	return true
}

// shouldBeRoot checks if a code should be considered a root node
func shouldBeRoot(code string) bool {
	// Main table divisions (0-9) should be roots
	if code == "0" || code == "1" || code == "2" || code == "3" ||
		code == "4" || code == "5" || code == "6" || code == "7" ||
		code == "8" || code == "9" {
		return true
	}

	// Top-level auxiliary signs should be roots
	if code == "+" || code == "/" || code == ":" || code == "::" ||
		code == "[]" || code == "*" || code == "A/Z" {
		return true
	}

	// Top-level auxiliary tables should be roots
	if code == "=..." || code == "(=...)" || code == "-0" {
		return true
	}

	// Single-digit language codes should be roots
	if strings.HasPrefix(code, "=") && len(code) == 2 && isNumeric(code[1:]) {
		return true
	}

	// Single-digit place codes should be roots
	if strings.HasPrefix(code, "(") && strings.HasSuffix(code, ")") {
		placeCode := code[1 : len(code)-1]
		if len(placeCode) == 1 && isNumeric(placeCode) {
			return true
		}
	}

	return false
}
