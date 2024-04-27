package parser

import (
	_ "embed"
	"encoding/json"
	"os"
	// "strings"
	"testing"

	// "github.com/andreyvit/diff"
)

//go:embed fixtures/test.hl
var fixture string

//go:embed fixtures/test.json
var fixtureResultSnapshot string

func TestLowlevelParse(t *testing.T) {
	parsed, err := Parse(fixture)
	if err != nil {
		t.Errorf("Error while parsing: %s", err)
	}

	contents, _ := json.MarshalIndent(parsed, "", "  ")
	os.WriteFile("fixtures/test.json", contents, 0644)
	// if strings.TrimSpace(fixtureResultSnapshot) == "update" {
	// 	os.WriteFile("fixtures/test.json", contents, 0644)
	// } else {
	// 	if string(contents) != fixtureResultSnapshot {
	// 		t.Errorf("Parsed result does not match snapshot:\n%v", diff.LineDiff(string(contents), fixtureResultSnapshot))
	// 	}
	// }
}
