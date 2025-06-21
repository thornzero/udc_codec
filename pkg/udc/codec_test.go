package udc

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestComposite(t *testing.T) {
	codec, err := LoadCodec("../../data/udc_full.yaml")
	if err != nil {
		t.Fatal(err)
	}
	nodes, err := codec.ParseComposite("621.3:681.5(075)")
	if err != nil {
		t.Fatal(err)
	}
	if len(nodes) != 3 {
		t.Errorf("expected 3 parts, got %d", len(nodes))
	}
}

func TestLoadCodec(t *testing.T) {
	codec, err := LoadCodec("../../data/udc_full.yaml")
	if err != nil {
		t.Fatal(err)
	}

	if codec == nil {
		t.Fatal("Expected codec to be non-nil")
	}

	if codec.flat == nil {
		t.Error("Expected flat map to be non-nil")
	}
}

func TestLookup(t *testing.T) {
	codec, err := LoadCodec("../../data/udc_full.yaml")
	if err != nil {
		t.Fatal(err)
	}

	// Test existing code
	title, found := codec.Lookup("0")
	if !found {
		t.Error("Expected to find code '0'")
	}
	if title == "" {
		t.Error("Expected non-empty title for code '0'")
	}

	// Test non-existing code
	_, found = codec.Lookup("NONEXISTENT")
	if found {
		t.Error("Expected not to find code 'NONEXISTENT'")
	}
}

func TestSearch(t *testing.T) {
	codec, err := LoadCodec("../../data/udc_full.yaml")
	if err != nil {
		t.Fatal(err)
	}

	// Test search for "computer"
	results := codec.Search("computer")
	if len(results) == 0 {
		t.Error("Expected to find results for 'computer'")
	}

	// Test search for non-existing term
	results = codec.Search("NONEXISTENTTERM")
	if len(results) > 0 {
		t.Error("Expected no results for 'NONEXISTENTTERM'")
	}
}

func TestChildren(t *testing.T) {
	codec, err := LoadCodec("../../data/udc_full.yaml")
	if err != nil {
		t.Fatal(err)
	}

	// Test children of existing code
	children, found := codec.Children("0")
	if !found {
		t.Error("Expected to find children for code '0'")
	}
	if len(children) == 0 {
		t.Error("Expected non-empty children for code '0'")
	}

	// Test children of non-existing code
	_, found = codec.Children("NONEXISTENT")
	if found {
		t.Error("Expected not to find children for code 'NONEXISTENT'")
	}
}

func TestAncestry(t *testing.T) {
	codec, err := LoadCodec("../../data/udc_full.yaml")
	if err != nil {
		t.Fatal(err)
	}

	// Test ancestry of existing code
	ancestry, found := codec.Ancestry("001")
	if !found {
		t.Error("Expected to find ancestry for code '001'")
	}
	if len(ancestry) == 0 {
		t.Error("Expected non-empty ancestry for code '001'")
	}

	// Test ancestry of non-existing code
	_, found = codec.Ancestry("NONEXISTENT")
	if found {
		t.Error("Expected not to find ancestry for code 'NONEXISTENT'")
	}
}

func TestValidate(t *testing.T) {
	codec, err := LoadCodec("../../data/udc_full.yaml")
	if err != nil {
		t.Fatal(err)
	}

	// Test valid composite code
	err = codec.Validate("621.3:681.5(075)")
	if err != nil {
		t.Errorf("Expected valid code '621.3:681.5(075)', got error: %v", err)
	}

	// Test invalid composite code
	err = codec.Validate("NONEXISTENT:ANOTHER")
	if err == nil {
		t.Error("Expected error for invalid code 'NONEXISTENT:ANOTHER'")
	}
}

func TestParseComposite(t *testing.T) {
	codec, err := LoadCodec("../../data/udc_full.yaml")
	if err != nil {
		t.Fatal(err)
	}

	// Test valid composite code
	nodes, err := codec.ParseComposite("621.3:681.5(075)")
	if err != nil {
		t.Fatalf("Expected valid code '621.3:681.5(075)', got error: %v", err)
	}
	if len(nodes) != 3 {
		t.Errorf("Expected 3 parts, got %d", len(nodes))
	}

	// Test invalid composite code
	_, err = codec.ParseComposite("NONEXISTENT:ANOTHER")
	if err == nil {
		t.Error("Expected error for invalid code 'NONEXISTENT:ANOTHER'")
	}
}

func TestLoadCodecNonExistentFile(t *testing.T) {
	// Test loading non-existent file
	_, err := LoadCodec("nonexistent_file.yaml")
	if err == nil {
		t.Error("Expected error when loading non-existent file")
	}
}

func TestLoadCodecWithOverlappingAddendum(t *testing.T) {
	// This test verifies that addendums with overlapping codes are rejected
	// We'll test this by creating a temporary addendum file with overlapping codes

	// Create a temporary addendum file with overlapping code
	tempDir := t.TempDir()
	addendumContent := `
- code: "0"
  title: "This should cause an error"
`
	addendumPath := filepath.Join(tempDir, "udc_addendum_overlap.yaml")
	err := os.WriteFile(addendumPath, []byte(addendumContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create temporary addendum file: %v", err)
	}

	// Try to load the codec with the overlapping addendum
	// This should fail because "0" is an existing UDC code
	_, err = LoadCodec("../../data/udc_full.yaml")
	if err == nil {
		t.Error("Expected error when loading addendum with overlapping code, but got none")
	} else if !strings.Contains(err.Error(), "overlapping code") {
		t.Errorf("Expected error about overlapping code, but got: %v", err)
	}
}

func TestLoadCodecWithValidAddendum(t *testing.T) {
	// This test verifies that valid addendums are loaded correctly
	// We'll test this by creating a temporary addendum file with valid codes

	// Create a temporary addendum file with valid codes
	tempDir := t.TempDir()
	addendumContent := `
- code: "999.1"
  title: "Test Local Classification"
  children:
    - code: "999.1.1"
      title: "Test Subclassification"
`
	addendumPath := filepath.Join(tempDir, "udc_addendum_valid.yaml")
	err := os.WriteFile(addendumPath, []byte(addendumContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create temporary addendum file: %v", err)
	}

	// Try to load the codec with the valid addendum
	// This should succeed because "999.1" is not an existing UDC code
	codec, err := LoadCodec("../../data/udc_full.yaml")
	if err != nil {
		t.Fatalf("Expected no error when loading valid addendum, but got: %v", err)
	}

	// Verify the addendum code is available
	title, found := codec.Lookup("999.1")
	if !found {
		t.Error("Expected to find addendum code '999.1'")
	}
	if title != "Test Local Classification" {
		t.Errorf("Expected title 'Test Local Classification' for code '999.1', got '%s'", title)
	}

	// Verify the addendum child code is available
	childTitle, found := codec.Lookup("999.1.1")
	if !found {
		t.Error("Expected to find addendum child code '999.1.1'")
	}
	if childTitle != "Test Subclassification" {
		t.Errorf("Expected title 'Test Subclassification' for code '999.1.1', got '%s'", childTitle)
	}
}

func TestAddendumManager(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "addendum_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create a mock UDC file so AddendumManager can validate codes
	udcContent := `
- code: "0"
  title: "Science and Knowledge"
- code: "1"
  title: "Philosophy. Psychology"
`
	udcFile := filepath.Join(tempDir, "udc_full.yaml")
	err = os.WriteFile(udcFile, []byte(udcContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create mock UDC file: %v", err)
	}

	am := NewAddendumManager(tempDir)

	// Test adding to default addendum
	t.Run("Add to default addendum", func(t *testing.T) {
		nodes := []*Node{
			{Code: "123.45", Title: "Test Classification"},
		}

		err := am.Add("", nodes)
		if err != nil {
			t.Fatalf("Failed to add to default addendum: %v", err)
		}

		// Check that file was created
		filepath := filepath.Join(tempDir, "udc_addendum_default.yaml")
		if _, err := os.Stat(filepath); os.IsNotExist(err) {
			t.Fatal("Default addendum file was not created")
		}

		// Verify content
		data, err := os.ReadFile(filepath)
		if err != nil {
			t.Fatal(err)
		}

		var loadedNodes []*Node
		if err := yaml.Unmarshal(data, &loadedNodes); err != nil {
			t.Fatal(err)
		}

		if len(loadedNodes) != 1 {
			t.Fatalf("Expected 1 node, got %d", len(loadedNodes))
		}

		if loadedNodes[0].Code != "123.45" || loadedNodes[0].Title != "Test Classification" {
			t.Fatalf("Unexpected node content: %+v", loadedNodes[0])
		}
	})

	// Test adding to specific addendum
	t.Run("Add to specific addendum", func(t *testing.T) {
		nodes := []*Node{
			{Code: "678.90", Title: "Another Classification"},
		}

		err := am.Add("custom", nodes)
		if err != nil {
			t.Fatalf("Failed to add to custom addendum: %v", err)
		}

		// Check that file was created
		filepath := filepath.Join(tempDir, "udc_addendum_custom.yaml")
		if _, err := os.Stat(filepath); os.IsNotExist(err) {
			t.Fatal("Custom addendum file was not created")
		}

		// Verify content
		data, err := os.ReadFile(filepath)
		if err != nil {
			t.Fatal(err)
		}

		var loadedNodes []*Node
		if err := yaml.Unmarshal(data, &loadedNodes); err != nil {
			t.Fatal(err)
		}

		if len(loadedNodes) != 1 {
			t.Fatalf("Expected 1 node, got %d", len(loadedNodes))
		}

		if loadedNodes[0].Code != "678.90" || loadedNodes[0].Title != "Another Classification" {
			t.Fatalf("Unexpected node content: %+v", loadedNodes[0])
		}
	})

	// Test adding to existing addendum
	t.Run("Add to existing addendum", func(t *testing.T) {
		// First add a node
		nodes1 := []*Node{
			{Code: "111.11", Title: "First Classification"},
		}
		err := am.Add("existing", nodes1)
		if err != nil {
			t.Fatalf("Failed to create addendum: %v", err)
		}

		// Then add another node to the same file
		nodes2 := []*Node{
			{Code: "222.22", Title: "Second Classification"},
		}
		err = am.Add("existing", nodes2)
		if err != nil {
			t.Fatalf("Failed to add to existing addendum: %v", err)
		}

		// Verify both nodes are present
		filepath := filepath.Join(tempDir, "udc_addendum_existing.yaml")
		data, err := os.ReadFile(filepath)
		if err != nil {
			t.Fatal(err)
		}

		var loadedNodes []*Node
		if err := yaml.Unmarshal(data, &loadedNodes); err != nil {
			t.Fatal(err)
		}

		if len(loadedNodes) != 2 {
			t.Fatalf("Expected 2 nodes, got %d", len(loadedNodes))
		}

		// Check both nodes are present (order may vary)
		codes := make(map[string]string)
		for _, node := range loadedNodes {
			codes[node.Code] = node.Title
		}

		if codes["111.11"] != "First Classification" {
			t.Fatalf("First classification not found or incorrect")
		}
		if codes["222.22"] != "Second Classification" {
			t.Fatalf("Second classification not found or incorrect")
		}
	})

	// Test validation - duplicate code
	t.Run("Validation - duplicate code", func(t *testing.T) {
		// Add a node first
		nodes1 := []*Node{
			{Code: "333.33", Title: "Original Classification"},
		}
		err := am.Add("validation", nodes1)
		if err != nil {
			t.Fatalf("Failed to create addendum: %v", err)
		}

		// Try to add a node with the same code
		nodes2 := []*Node{
			{Code: "333.33", Title: "Duplicate Classification"},
		}
		err = am.Add("validation", nodes2)
		if err == nil {
			t.Fatal("Expected error for duplicate code, got none")
		}

		if !strings.Contains(err.Error(), "already exists") {
			t.Fatalf("Expected duplicate code error, got: %v", err)
		}
	})

	// Test listing addendums
	t.Run("List addendums", func(t *testing.T) {
		addendums, err := am.ListAddendums()
		if err != nil {
			t.Fatalf("Failed to list addendums: %v", err)
		}

		expected := []string{
			"udc_addendum_default.yaml",
			"udc_addendum_custom.yaml",
			"udc_addendum_existing.yaml",
			"udc_addendum_validation.yaml",
		}

		if len(addendums) != len(expected) {
			t.Fatalf("Expected %d addendums, got %d", len(expected), len(addendums))
		}

		for _, expectedFile := range expected {
			found := false
			for _, addendum := range addendums {
				if addendum == expectedFile {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("Expected addendum %s not found", expectedFile)
			}
		}
	})

	// Test deleting addendum
	t.Run("Delete addendum", func(t *testing.T) {
		err := am.DeleteAddendum("udc_addendum_custom.yaml")
		if err != nil {
			t.Fatalf("Failed to delete addendum: %v", err)
		}

		// Verify file was deleted
		filepath := filepath.Join(tempDir, "udc_addendum_custom.yaml")
		if _, err := os.Stat(filepath); !os.IsNotExist(err) {
			t.Fatal("Addendum file was not deleted")
		}

		// Verify it's not in the list anymore
		addendums, err := am.ListAddendums()
		if err != nil {
			t.Fatalf("Failed to list addendums: %v", err)
		}

		for _, addendum := range addendums {
			if addendum == "udc_addendum_custom.yaml" {
				t.Fatal("Deleted addendum still appears in list")
			}
		}
	})
}

func TestAddendumManagerValidation(t *testing.T) {
	tempDir := t.TempDir()

	// Create a mock UDC file
	udcContent := `
- code: "0"
  title: "Science and Knowledge"
`
	udcFile := filepath.Join(tempDir, "udc_full.yaml")
	err := os.WriteFile(udcFile, []byte(udcContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create mock UDC file: %v", err)
	}

	am := NewAddendumManager(tempDir)

	// Test invalid filename format (should be auto-corrected now)
	invalidNodes := []*Node{{Code: "999.1", Title: "Test"}}

	err = am.Add("invalid_name", invalidNodes)
	if err != nil {
		t.Errorf("Expected no error for invalid filename (should be auto-corrected), but got: %v", err)
	}

	// Test creating addendum with child that overlaps
	overlappingChildNodes := []*Node{
		{
			Code:  "999.1",
			Title: "Valid Parent",
			Children: []*Node{
				{
					Code:  "0",
					Title: "Invalid Child",
				},
			},
		},
	}

	err = am.Add("invalid_child", overlappingChildNodes)
	if err == nil {
		t.Error("Expected error when creating addendum with overlapping child code")
	} else if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("Expected error about existing code, but got: %v", err)
	}
}
