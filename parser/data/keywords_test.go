package parser_data

import "testing"

func TestFindKeyword(t *testing.T) {
	k, found := FindKeyword("submap")
	if !found {
		t.Fatal("submap not found")
	}
	if k.Name != "submap" {
		t.Fatalf("unexpected name: %q", k.Name)
	}

	k, found = FindKeyword("binde")
	if !found {
		t.Fatal("binde not found")
	}
	if k.Name != "bind" {
		t.Fatalf("unexpected name: %q", k.Name)
	}

}
