package assettag

import "testing"

func TestParser(t *testing.T) {
	tag, err := ParseTag("POL-LT1001")
	if err != nil {
		t.Fatal(err)
	}
	if tag.FunctionCode != "LT" {
		t.Errorf("expected LT, got %s", tag.FunctionCode)
	}
}
