package parser

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestHighLevelParse(t *testing.T) {
	parsed, err := Parse(fixture)
	if err != nil {
		t.Errorf("Error while parsing: %s", err)
	}

	config, err := parsed.Decode()
	if err != nil {
		t.Errorf("Error while decoding: %s", err)
		return
	} else {
		spew.Dump(config)
	}

	t.Errorf("TestHighLevelParse not implemented")
}
