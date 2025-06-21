package udc

import (
	"testing"
)

func TestConvertRawToModel(t *testing.T) {
	// Create a simple hierarchy of RawNodes
	rawNodes := []*RawNode{
		{
			ID:    "1",
			Code:  "TOP",
			Title: "UDC Summary Root",
			Children: []*RawNode{
				{
					ID:    "2",
					Code:  "0",
					Title: "Science and Knowledge",
					Children: []*RawNode{
						{
							ID:       "3",
							Code:     "00",
							Title:    "Prolegomena",
							Children: []*RawNode{},
						},
					},
				},
			},
		},
	}

	// Convert to UDCNode
	udcNodes := ConvertRawToModel(rawNodes)

	// Verify conversion
	if len(udcNodes) != 1 {
		t.Fatalf("Expected 1 root node, got %d", len(udcNodes))
	}

	root := udcNodes[0]
	if root.Code != "TOP" {
		t.Errorf("Expected root code to be 'TOP', got '%s'", root.Code)
	}
	if root.Title != "UDC Summary Root" {
		t.Errorf("Expected root title to be 'UDC Summary Root', got '%s'", root.Title)
	}

	if len(root.Children) != 1 {
		t.Fatalf("Expected 1 child, got %d", len(root.Children))
	}

	child := root.Children[0]
	if child.Code != "0" {
		t.Errorf("Expected child code to be '0', got '%s'", child.Code)
	}
	if child.Title != "Science and Knowledge" {
		t.Errorf("Expected child title to be 'Science and Knowledge', got '%s'", child.Title)
	}

	if len(child.Children) != 1 {
		t.Fatalf("Expected 1 grandchild, got %d", len(child.Children))
	}

	grandchild := child.Children[0]
	if grandchild.Code != "00" {
		t.Errorf("Expected grandchild code to be '00', got '%s'", grandchild.Code)
	}
	if grandchild.Title != "Prolegomena" {
		t.Errorf("Expected grandchild title to be 'Prolegomena', got '%s'", grandchild.Title)
	}

	if len(grandchild.Children) != 0 {
		t.Errorf("Expected 0 great-grandchildren, got %d", len(grandchild.Children))
	}
}

func TestConvertRawToModelEmpty(t *testing.T) {
	// Test with empty input
	rawNodes := []*RawNode{}
	udcNodes := ConvertRawToModel(rawNodes)

	if len(udcNodes) != 0 {
		t.Errorf("Expected 0 nodes for empty input, got %d", len(udcNodes))
	}
}

func TestConvertRawToModelSingleNode(t *testing.T) {
	// Test with single node without children
	rawNodes := []*RawNode{
		{
			ID:       "1",
			Code:     "TEST",
			Title:    "Test Node",
			Children: []*RawNode{},
		},
	}

	udcNodes := ConvertRawToModel(rawNodes)

	if len(udcNodes) != 1 {
		t.Fatalf("Expected 1 node, got %d", len(udcNodes))
	}

	node := udcNodes[0]
	if node.Code != "TEST" {
		t.Errorf("Expected code to be 'TEST', got '%s'", node.Code)
	}
	if node.Title != "Test Node" {
		t.Errorf("Expected title to be 'Test Node', got '%s'", node.Title)
	}
	if len(node.Children) != 0 {
		t.Errorf("Expected 0 children, got %d", len(node.Children))
	}
}

func TestConvertRawToModelMultipleRoots(t *testing.T) {
	// Test with multiple root nodes
	rawNodes := []*RawNode{
		{
			ID:       "1",
			Code:     "ROOT1",
			Title:    "Root 1",
			Children: []*RawNode{},
		},
		{
			ID:       "2",
			Code:     "ROOT2",
			Title:    "Root 2",
			Children: []*RawNode{},
		},
	}

	udcNodes := ConvertRawToModel(rawNodes)

	if len(udcNodes) != 2 {
		t.Fatalf("Expected 2 root nodes, got %d", len(udcNodes))
	}

	if udcNodes[0].Code != "ROOT1" {
		t.Errorf("Expected first root code to be 'ROOT1', got '%s'", udcNodes[0].Code)
	}
	if udcNodes[1].Code != "ROOT2" {
		t.Errorf("Expected second root code to be 'ROOT2', got '%s'", udcNodes[1].Code)
	}
}
