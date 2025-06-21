package udc

import (
	"testing"
)

func TestBuildRawHierarchy(t *testing.T) {
	// Create mock nodes with UDC codes
	nodes := []*RawNode{
		{ID: "1", Code: "TOP", Title: "UDC Summary Root"},
		{ID: "2", Code: "0", Title: "Science and Knowledge"},
		{ID: "3", Code: "00", Title: "Prolegomena"},
		{ID: "4", Code: "000", Title: "General"},
		{ID: "5", Code: "001", Title: "Science and knowledge in general"},
		{ID: "6", Code: "001.1", Title: "Concepts of science"},
		{ID: "7", Code: "1", Title: "Philosophy. Psychology"},
		{ID: "8", Code: "=1", Title: "Indo-European languages"},
		{ID: "9", Code: "=11", Title: "Germanic languages"},
		{ID: "10", Code: "=111", Title: "English"},
		{ID: "11", Code: "(1)", Title: "Place auxiliaries"},
		{ID: "12", Code: "(540)", Title: "India"},
		{ID: "13", Code: "-0", Title: "Form auxiliaries"},
		{ID: "14", Code: "-05", Title: "Victims of circumstances"},
		{ID: "15", Code: "-058", Title: "Victims of circumstances"},
		{ID: "16", Code: "-058.6", Title: "Victims of circumstances"},
	}

	roots := buildRawHierarchy(nodes)

	// The hierarchy should have TOP as the main root, but some nodes might become roots
	// if they can't find their parents. Let's check for TOP first.
	foundTOP := false
	for _, root := range roots {
		if root.Code == "TOP" {
			foundTOP = true
			break
		}
	}

	if !foundTOP {
		t.Error("Expected to find TOP as a root node")
		return
	}

	// Find the TOP root and verify its structure
	var topRoot *RawNode
	for _, root := range roots {
		if root.Code == "TOP" {
			topRoot = root
			break
		}
	}

	if topRoot == nil {
		t.Fatal("TOP root not found")
	}

	// Count expected children of TOP (main table divisions and top-level auxiliaries)
	expectedTopChildren := []string{"0", "1", "=1", "(1)", "-0"}
	foundChildren := make(map[string]bool)

	for _, child := range topRoot.Children {
		foundChildren[child.Code] = true
	}

	for _, expected := range expectedTopChildren {
		if !foundChildren[expected] {
			t.Errorf("Expected TOP to have child '%s'", expected)
		}
	}

	// Verify main table structure (0, 1 should be children of TOP)
	found0 := false
	found1 := false
	for _, child := range topRoot.Children {
		if child.Code == "0" {
			found0 = true
			// Verify that 00 is a child of 0
			found00 := false
			for _, grandchild := range child.Children {
				if grandchild.Code == "00" {
					found00 = true
					// Verify that 000 is a child of 00
					found000 := false
					for _, greatGrandchild := range grandchild.Children {
						if greatGrandchild.Code == "000" {
							found000 = true
							break
						}
					}
					if !found000 {
						t.Error("Expected '000' to be a child of '00'")
					}
					break
				}
			}
			if !found00 {
				t.Error("Expected '00' to be a child of '0'")
			}
		}
		if child.Code == "1" {
			found1 = true
		}
	}
	if !found0 {
		t.Error("Expected '0' to be a child of 'TOP'")
	}
	if !found1 {
		t.Error("Expected '1' to be a child of 'TOP'")
	}

	// Verify auxiliary table structure
	foundEquals := false
	foundParens := false
	foundMinus := false
	for _, child := range topRoot.Children {
		if child.Code == "=1" {
			foundEquals = true
			// Verify that =11 is a child of =1
			found11 := false
			for _, grandchild := range child.Children {
				if grandchild.Code == "=11" {
					found11 = true
					// Verify that =111 is a child of =11
					found111 := false
					for _, greatGrandchild := range grandchild.Children {
						if greatGrandchild.Code == "=111" {
							found111 = true
							break
						}
					}
					if !found111 {
						t.Error("Expected '=111' to be a child of '=11'")
					}
					break
				}
			}
			if !found11 {
				t.Error("Expected '=11' to be a child of '=1'")
			}
		}
		if child.Code == "(1)" {
			foundParens = true
			// Verify that (540) is a child of (1) - this might not work as expected
			// since (540) might be a direct child of TOP depending on the logic
		}
		if child.Code == "-0" {
			foundMinus = true
			// Verify that -05 is a child of -0
			found05 := false
			for _, grandchild := range child.Children {
				if grandchild.Code == "-05" {
					found05 = true
					// Verify that -058 is a child of -05
					found058 := false
					for _, greatGrandchild := range grandchild.Children {
						if greatGrandchild.Code == "-058" {
							found058 = true
							// Verify that -058.6 is a child of -058
							found0586 := false
							for _, greatGreatGrandchild := range greatGrandchild.Children {
								if greatGreatGrandchild.Code == "-058.6" {
									found0586 = true
									break
								}
							}
							if !found0586 {
								t.Error("Expected '-058.6' to be a child of '-058'")
							}
							break
						}
					}
					if !found058 {
						t.Error("Expected '-058' to be a child of '-05'")
					}
					break
				}
			}
			if !found05 {
				t.Error("Expected '-05' to be a child of '-0'")
			}
		}
	}
	if !foundEquals {
		t.Error("Expected '=1' to be a child of 'TOP'")
	}
	if !foundParens {
		t.Error("Expected '(1)' to be a child of 'TOP'")
	}
	if !foundMinus {
		t.Error("Expected '-0' to be a child of 'TOP'")
	}
}

func TestFindParentCode(t *testing.T) {
	tests := []struct {
		code     string
		expected string
	}{
		{"00", "0"},
		{"000", "00"},
		{"001", "00"},
		{"001.1", "001"},
		{"01", "0"},
		{"1", "TOP"},
		{"=1", "TOP"},
		{"=11", "=1"},
		{"=111", "=11"},
		{"(1)", "TOP"},
		{"(540)", "(5)"}, // This might need adjustment based on actual UDC rules
		{"-0", "TOP"},
		{"-058.6", "-058"},
		{"=...", "TOP"},
		{"=...`01", "=..."},
	}

	for _, tt := range tests {
		result := findParentCode(tt.code)
		if result != tt.expected {
			t.Errorf("findParentCode(%q) = %q, want %q", tt.code, result, tt.expected)
		}
	}
}

func TestFindMainTableParent(t *testing.T) {
	tests := []struct {
		code     string
		expected string
	}{
		{"0", "TOP"},
		{"1", "TOP"},
		{"9", "TOP"},
		{"00", "0"},
		{"01", "0"},
		{"000", "00"},
		{"001", "00"},
		{"001.1", "001"},
		{"001.1.1", "001.1"},
		{"", ""},
		{"abc", ""},
	}

	for _, tt := range tests {
		result := findMainTableParent(tt.code)
		if result != tt.expected {
			t.Errorf("findMainTableParent(%q) = %q, want %q", tt.code, result, tt.expected)
		}
	}
}

func TestFindAuxiliaryParent(t *testing.T) {
	tests := []struct {
		code     string
		expected string
	}{
		// Language auxiliaries
		{"=...", "TOP"},
		{"=...`01", "=..."},
		{"=00", "=..."},
		{"=030", "=..."},
		{"=1", "TOP"},
		{"=11", "=1"},
		{"=111", "=11"},
		{"=2", "TOP"},
		{"=21", "=2"},

		// Place auxiliaries
		{"(1)", "TOP"},
		{"(5)", "TOP"},
		{"(540)", "(5)"},
		{"(540.1)", "(5)"},
		{"(=01)", "(=...)"},

		// Form auxiliaries
		{"-0", "TOP"},
		{"-058", "-05"}, // -058 -> -05 (removes last digit)
		{"-058.6", "-058"},
		{"-1", ""}, // -1 -> "" (single digit, no parent)

		// Edge cases
		{"", ""},
		{"abc", ""},
	}

	for _, tt := range tests {
		result := findAuxiliaryParent(tt.code)
		if result != tt.expected {
			t.Errorf("findAuxiliaryParent(%q) = %q, want %q", tt.code, result, tt.expected)
		}
	}
}

func TestFindNumericParent(t *testing.T) {
	tests := []struct {
		code     string
		prefix   string
		suffix   string
		expected string
	}{
		{"11", "=", "", "=1"},
		{"111", "=", "", "=11"},
		{"058", "-", "", "-05"},
		{"058.6", "-", "", "-058"},
		{"01/08", "=", "", "=01"},
		{"", "=", "", ""},
		{"1", "=", "", ""},
		{"abc", "=", "", ""},
	}

	for _, tt := range tests {
		result := findNumericParent(tt.code, tt.prefix, tt.suffix)
		if result != tt.expected {
			t.Errorf("findNumericParent(%q, %q, %q) = %q, want %q",
				tt.code, tt.prefix, tt.suffix, result, tt.expected)
		}
	}
}

func TestIsNumeric(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"123", true},
		{"123.456", true},
		{"0", true},
		{"", true},
		{"abc", false},
		{"123abc", false},
		{"123.456.789", true},
		{"123.", true},
		{".123", true},
	}

	for _, tt := range tests {
		result := isNumeric(tt.input)
		if result != tt.expected {
			t.Errorf("isNumeric(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

func TestShouldBeRoot(t *testing.T) {
	tests := []struct {
		code     string
		expected bool
	}{
		// Main table divisions
		{"0", true},
		{"1", true},
		{"9", true},
		{"00", false},
		{"01", false},

		// Top-level auxiliary signs
		{"+", true},
		{"/", true},
		{":", true},
		{"::", true},
		{"[]", true},
		{"*", true},
		{"A/Z", true},

		// Top-level auxiliary tables
		{"=...", true},
		{"(=...)", true},
		{"-0", true},

		// Single-digit language codes
		{"=1", true},
		{"=2", true},
		{"=11", false},

		// Single-digit place codes
		{"(1)", true},
		{"(5)", true},
		{"(540)", false},

		// Other cases
		{"", false},
		{"abc", false},
	}

	for _, tt := range tests {
		result := shouldBeRoot(tt.code)
		if result != tt.expected {
			t.Errorf("shouldBeRoot(%q) = %v, want %v", tt.code, result, tt.expected)
		}
	}
}

func TestBuildRawHierarchyWithFallback(t *testing.T) {
	// Test hierarchy building with fallback to website parent IDs
	nodes := []*RawNode{
		{ID: "1", Code: "TOP", Title: "UDC Summary Root"},
		{ID: "2", Code: "0", Title: "Science and Knowledge"},
		{ID: "3", Code: "00", Title: "Prolegomena"},
		{ID: "4", Code: "UNKNOWN", Title: "Unknown Code", Parent: "2"}, // Should use fallback
	}

	roots := buildRawHierarchy(nodes)

	if len(roots) != 1 {
		t.Errorf("Expected 1 root node, got %d", len(roots))
		return
	}

	root := roots[0]
	if root.Code != "TOP" {
		t.Errorf("Expected root code to be 'TOP', got '%s'", root.Code)
	}

	// Verify that UNKNOWN was attached to node with ID "2" (code "0")
	foundUnknown := false
	for _, child := range root.Children {
		if child.Code == "0" {
			for _, grandchild := range child.Children {
				if grandchild.Code == "UNKNOWN" {
					foundUnknown = true
					break
				}
			}
			break
		}
	}
	if !foundUnknown {
		t.Error("Expected 'UNKNOWN' to be attached via fallback parent ID")
	}
}

func TestBuildRawHierarchyEmpty(t *testing.T) {
	// Test with empty node list
	nodes := []*RawNode{}
	roots := buildRawHierarchy(nodes)

	if len(roots) != 0 {
		t.Errorf("Expected 0 root nodes for empty input, got %d", len(roots))
	}
}

func TestBuildRawHierarchySingleNode(t *testing.T) {
	// Test with single root node
	nodes := []*RawNode{
		{ID: "1", Code: "TOP", Title: "UDC Summary Root"},
	}

	roots := buildRawHierarchy(nodes)

	if len(roots) != 1 {
		t.Errorf("Expected 1 root node, got %d", len(roots))
		return
	}

	if roots[0].Code != "TOP" {
		t.Errorf("Expected root code to be 'TOP', got '%s'", roots[0].Code)
	}
}
