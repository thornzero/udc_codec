package udc

import "testing"

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
