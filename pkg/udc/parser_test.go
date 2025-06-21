package udc

import (
	"regexp"
	"testing"
)

func TestParseRawHTML(t *testing.T) {
	// Test HTML with valid d.add() calls - matching the expected regex pattern exactly
	html := `
		<div id="classtree">
			<script>
				d.add(1, -1, 'TOP', '<span class="nodetag">TOP</span>&nbsp;&nbsp;UDC Summary Root', 'TOP');
				d.add(2, 1, '0', '<span class="nodetag">0</span>&nbsp;&nbsp;Science and Knowledge', '0');
				d.add(3, 2, '00', '<span class="nodetag">00</span>&nbsp;&nbsp;Prolegomena', '00');
				d.add(4, 3, '001', '<span class="nodetag">001</span>&nbsp;&nbsp;Science and knowledge in general', '001');
			</script>
		</div>
	`

	nodes := parseRawHTML(html)

	if len(nodes) != 4 {
		t.Errorf("Expected 4 nodes, got %d", len(nodes))
		return
	}

	// Test first node
	if nodes[0].ID != "1" {
		t.Errorf("Expected ID '1', got '%s'", nodes[0].ID)
	}
	if nodes[0].Code != "TOP" {
		t.Errorf("Expected code 'TOP', got '%s'", nodes[0].Code)
	}
	if nodes[0].Title != "UDC Summary Root" {
		t.Errorf("Expected title 'UDC Summary Root', got '%s'", nodes[0].Title)
	}
	if nodes[0].Parent != "-1" {
		t.Errorf("Expected parent '-1', got '%s'", nodes[0].Parent)
	}

	// Test second node
	if nodes[1].ID != "2" {
		t.Errorf("Expected ID '2', got '%s'", nodes[1].ID)
	}
	if nodes[1].Code != "0" {
		t.Errorf("Expected code '0', got '%s'", nodes[1].Code)
	}
	if nodes[1].Title != "Science and Knowledge" {
		t.Errorf("Expected title 'Science and Knowledge', got '%s'", nodes[1].Title)
	}
	if nodes[1].Parent != "1" {
		t.Errorf("Expected parent '1', got '%s'", nodes[1].Parent)
	}
}

func TestParseRawHTMLWithInvalidCodes(t *testing.T) {
	// Test HTML with invalid codes that should be filtered out
	html := `
		<div id="classtree">
			<script>
				d.add(1, -1, 'TOP', '<span class="nodetag">TOP</span>&nbsp;&nbsp;UDC Summary Root', 'TOP');
				d.add(2, 1, '-', '<span class="nodetag">-</span>&nbsp;&nbsp;Invalid Code', '-');
				d.add(3, 1, '--', '<span class="nodetag">--</span>&nbsp;&nbsp;Invalid Code', '--');
				d.add(4, 1, '---', '<span class="nodetag">---</span>&nbsp;&nbsp;Invalid Code', '---');
				d.add(5, 1, '----', '<span class="nodetag">----</span>&nbsp;&nbsp;Invalid Code', '----');
				d.add(6, 1, '0', '<span class="nodetag">0</span>&nbsp;&nbsp;Valid Code', '0');
			</script>
		</div>
	`

	nodes := parseRawHTML(html)

	// Should only have 2 nodes (TOP and 0), invalid codes should be filtered out
	if len(nodes) != 2 {
		t.Errorf("Expected 2 nodes (invalid codes filtered), got %d", len(nodes))
		return
	}

	// Check that only valid codes remain
	validCodes := []string{"TOP", "0"}
	for i, expectedCode := range validCodes {
		if nodes[i].Code != expectedCode {
			t.Errorf("Expected code '%s' at position %d, got '%s'", expectedCode, i, nodes[i].Code)
		}
	}
}

func TestParseRawHTMLEmpty(t *testing.T) {
	// Test with empty HTML
	html := ""
	nodes := parseRawHTML(html)

	if len(nodes) != 0 {
		t.Errorf("Expected 0 nodes for empty HTML, got %d", len(nodes))
	}
}

func TestParseRawHTMLNoMatches(t *testing.T) {
	// Test HTML without any d.add() calls
	html := `
		<div id="classtree">
			<p>Some text without d.add() calls</p>
		</div>
	`
	nodes := parseRawHTML(html)

	if len(nodes) != 0 {
		t.Errorf("Expected 0 nodes for HTML without d.add() calls, got %d", len(nodes))
	}
}

func TestParseRawHTMLDebug(t *testing.T) {
	// Test HTML with valid d.add() calls - matching the expected regex pattern exactly
	html := `
		<div id="classtree">
			<script>
				d.add(1, -1, 'TOP', '<span class="nodetag">TOP</span>&nbsp;&nbsp;UDC Summary Root', 'TOP');
				d.add(2, 1, '0', '<span class="nodetag">0</span>&nbsp;&nbsp;Science and Knowledge', '0');
			</script>
		</div>
	`

	// Test the regex pattern directly - using the same pattern as the implementation
	re := regexp.MustCompile(`d\.add\(\s*(\d+),\s*(-\d|\d+),\s*\'(.*?)\'\,\s*\'[^\']*?&nbsp;&nbsp;(.*?)\'`)
	matches := re.FindAllStringSubmatch(html, -1)

	t.Logf("Found %d matches", len(matches))
	for i, match := range matches {
		t.Logf("Match %d: %v", i, match)
	}

	if len(matches) != 2 {
		t.Errorf("Expected 2 matches, got %d", len(matches))
	}
}
