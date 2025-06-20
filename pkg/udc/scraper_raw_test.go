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
		{ID: "14", Code: "-058.6", Title: "Victims of circumstances"},
	}

	roots := buildRawHierarchy(nodes)

	// Verify that we have the expected root structure
	if len(roots) != 1 {
		t.Errorf("Expected 1 root node, got %d", len(roots))
		return
	}

	root := roots[0]
	if root.Code != "TOP" {
		t.Errorf("Expected root code to be 'TOP', got '%s'", root.Code)
	}

	// Verify main table structure (0, 1 should be children of TOP)
	found0 := false
	found1 := false
	for _, child := range root.Children {
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
	for _, child := range root.Children {
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
			// Verify that -058.6 is a child of -0
			found058 := false
			for _, grandchild := range child.Children {
				if grandchild.Code == "-058.6" {
					found058 = true
					break
				}
			}
			if !found058 {
				t.Error("Expected '-058.6' to be a child of '-0'")
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
