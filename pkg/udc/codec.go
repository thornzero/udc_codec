package udc

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Codec struct {
	flat map[string]*Node
}

type Node struct {
	Code     string  `yaml:"code"`
	Title    string  `yaml:"title"`
	Children []*Node `yaml:"children,omitempty"`
}

// LoadCodec loads the UDC codec from udc_full.yaml and merges any local addendums
func LoadCodec(filename string) (*Codec, error) {
	// Load the main UDC data
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", filename, err)
	}

	var nodes []*Node
	if err := yaml.Unmarshal(data, &nodes); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", filename, err)
	}

	// Load and merge local addendums
	addendumNodes, err := loadAddendums(filepath.Dir(filename))
	if err != nil {
		return nil, fmt.Errorf("failed to load addendums: %w", err)
	}

	// Merge addendum nodes with main nodes
	nodes, err = mergeNodes(nodes, addendumNodes)
	if err != nil {
		return nil, err
	}

	// Build flat map
	flat := make(map[string]*Node)
	buildFlatMap(nodes, flat)

	return &Codec{flat: flat}, nil
}

// loadAddendums loads all addendum files from the data directory
func loadAddendums(dataDir string) ([]*Node, error) {
	var allAddendumNodes []*Node

	// Look for files matching the pattern "udc_addendum_*.yaml"
	err := filepath.WalkDir(dataDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-yaml files
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".yaml") {
			return nil
		}

		// Check if it's an addendum file
		if !strings.HasPrefix(d.Name(), "udc_addendum_") {
			return nil
		}

		// Load the addendum file
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read addendum %s: %w", path, err)
		}

		var addendumNodes []*Node
		if err := yaml.Unmarshal(data, &addendumNodes); err != nil {
			return fmt.Errorf("failed to parse addendum %s: %w", path, err)
		}

		allAddendumNodes = append(allAddendumNodes, addendumNodes...)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return allAddendumNodes, nil
}

// mergeNodes merges addendum nodes with main nodes, rejecting overlaps
func mergeNodes(mainNodes, addendumNodes []*Node) ([]*Node, error) {
	// Create a map of main nodes by code for easy lookup
	mainMap := make(map[string]*Node)
	for _, node := range mainNodes {
		mainMap[node.Code] = node
	}

	// Check for overlaps and collect new nodes
	var newNodes []*Node
	for _, addendumNode := range addendumNodes {
		if _, exists := mainMap[addendumNode.Code]; exists {
			return nil, fmt.Errorf("addendum contains overlapping code: %s (addendums cannot override existing UDC codes)", addendumNode.Code)
		}
		newNodes = append(newNodes, addendumNode)
	}

	// Return main nodes plus new nodes (no merging of existing codes)
	return append(mainNodes, newNodes...), nil
}

// buildFlatMap builds a flat map of all nodes by their code
func buildFlatMap(nodes []*Node, flat map[string]*Node) {
	for _, node := range nodes {
		flat[node.Code] = node
		if node.Children != nil {
			buildFlatMap(node.Children, flat)
		}
	}
}

func (c *Codec) Lookup(code string) (string, bool) {
	node, ok := c.flat[code]
	if !ok {
		return "", false
	}
	return node.Title, true
}

func (c *Codec) ParseComposite(code string) ([]*Node, error) {
	parts, err := parseComposite(code)
	if err != nil {
		return nil, err
	}
	var result []*Node
	for _, p := range parts {
		node, ok := c.flat[p]
		if !ok {
			return nil, fmt.Errorf("unknown code part: %s", p)
		}
		result = append(result, node)
	}
	return result, nil
}

func (c *Codec) Search(term string) []*Node {
	return search(c.flat, term)
}

func (c *Codec) Children(code string) ([]*Node, bool) {
	node, ok := c.flat[code]
	if !ok {
		return nil, false
	}
	return node.Children, true
}

func (c *Codec) Ancestry(code string) ([]*Node, bool) {
	// Since we don't have parent pointers in the flat map,
	// we need to reconstruct the ancestry by searching for parent codes
	var path []*Node
	currentCode := code
	for currentCode != "" {
		if currentNode, exists := c.flat[currentCode]; exists {
			path = append([]*Node{currentNode}, path...)
			// Find parent code by removing last character or dot-separated part
			currentCode = findParentCode(currentCode)
		} else {
			break
		}
	}
	return path, true
}

func (c *Codec) Validate(code string) error {
	return validateComposite(code, c.flat)
}

// AddendumManager provides functions for managing addendum files
type AddendumManager struct {
	dataDir string
}

// NewAddendumManager creates a new addendum manager for the given data directory
func NewAddendumManager(dataDir string) *AddendumManager {
	return &AddendumManager{dataDir: dataDir}
}

// Add adds nodes to an addendum file, creating it if it doesn't exist
// If filename is empty, uses the default addendum file
func (am *AddendumManager) Add(filename string, nodes []*Node) error {
	// Use default filename if none provided
	if filename == "" {
		filename = "udc_addendum_default.yaml"
	}

	// Ensure filename has correct format
	if !strings.HasPrefix(filename, "udc_addendum_") {
		filename = "udc_addendum_" + filename
	}
	if !strings.HasSuffix(filename, ".yaml") {
		filename = filename + ".yaml"
	}

	filepath := filepath.Join(am.dataDir, filename)

	// Load existing addendum if it exists
	var existingNodes []*Node
	if _, err := os.Stat(filepath); err == nil {
		data, err := os.ReadFile(filepath)
		if err != nil {
			return fmt.Errorf("failed to read existing addendum: %w", err)
		}
		if err := yaml.Unmarshal(data, &existingNodes); err != nil {
			return fmt.Errorf("failed to parse existing addendum: %w", err)
		}
	}

	// Get all existing codes (UDC + existing addendums)
	existingCodes, err := am.getExistingCodes()
	if err != nil {
		return fmt.Errorf("failed to get existing codes: %w", err)
	}

	// Add existing addendum codes to the set
	for _, node := range existingNodes {
		existingCodes[node.Code] = true
		if node.Children != nil {
			am.collectChildCodes(node.Children, existingCodes)
		}
	}

	// Validate new nodes
	for _, node := range nodes {
		if err := am.validateNode(node, existingCodes); err != nil {
			return err
		}
	}

	// Merge nodes
	allNodes := append(existingNodes, nodes...)

	// Marshal and save
	data, err := yaml.Marshal(allNodes)
	if err != nil {
		return fmt.Errorf("failed to marshal addendum data: %w", err)
	}

	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write addendum file: %w", err)
	}

	return nil
}

// ListAddendums returns a list of all addendum files
func (am *AddendumManager) ListAddendums() ([]string, error) {
	files, err := os.ReadDir(am.dataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read data directory: %w", err)
	}

	var addendums []string
	for _, file := range files {
		if !file.IsDir() && strings.HasPrefix(file.Name(), "udc_addendum_") && strings.HasSuffix(file.Name(), ".yaml") {
			addendums = append(addendums, file.Name())
		}
	}

	return addendums, nil
}

// DeleteAddendum deletes an addendum file
func (am *AddendumManager) DeleteAddendum(filename string) error {
	// Ensure filename has correct format
	if !strings.HasPrefix(filename, "udc_addendum_") {
		filename = "udc_addendum_" + filename
	}
	if !strings.HasSuffix(filename, ".yaml") {
		filename = filename + ".yaml"
	}

	filepath := filepath.Join(am.dataDir, filename)
	return os.Remove(filepath)
}

// getExistingCodes gets all existing UDC codes from the main file
func (am *AddendumManager) getExistingCodes() (map[string]bool, error) {
	udcFile := filepath.Join(am.dataDir, "udc_full.yaml")

	data, err := os.ReadFile(udcFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read UDC file: %w", err)
	}

	var nodes []*Node
	if err := yaml.Unmarshal(data, &nodes); err != nil {
		return nil, fmt.Errorf("failed to parse UDC file: %w", err)
	}

	codes := make(map[string]bool)
	am.collectChildCodes(nodes, codes)
	return codes, nil
}

// collectChildCodes recursively collects all codes from nodes and their children
func (am *AddendumManager) collectChildCodes(nodes []*Node, codes map[string]bool) {
	for _, node := range nodes {
		codes[node.Code] = true
		if node.Children != nil {
			am.collectChildCodes(node.Children, codes)
		}
	}
}

// validateNode validates that a node and its children don't overlap with existing codes
func (am *AddendumManager) validateNode(node *Node, existingCodes map[string]bool) error {
	if existingCodes[node.Code] {
		return fmt.Errorf("code '%s' already exists in UDC classification", node.Code)
	}

	if node.Children != nil {
		for _, child := range node.Children {
			if err := am.validateNode(child, existingCodes); err != nil {
				return err
			}
		}
	}

	return nil
}
