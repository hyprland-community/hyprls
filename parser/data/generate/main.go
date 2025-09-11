package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/MakeNowJust/heredoc/v2"
	. "github.com/hyprland-community/hyprls/parser/data"
)

func main() {
	rootSections := make([]SectionDefinition, 0)
	for _, section := range Sections {
		if len(section.Path) == 1 {
			rootSections = append(rootSections, section)
		}
	}

	jsoned, _ := json.Marshal(struct {
		Sections []SectionDefinition `json:"sections"`
		Keywords []KeywordDefinition `json:"keywords"`
	}{
		Sections: rootSections,
		Keywords: Keywords,
	})
	os.WriteFile(os.Args[1], jsoned, 0644)

	fmt.Println(heredoc.Doc(`package parser

		import "image/color"

		type Configuration struct {
			CustomVariables map[string]string
	`))

	for _, section := range rootSections {
		fmt.Printf("\t%s %s `json:\"%s\"`\n", section.Name(), section.TypeName(), section.JSONName())
	}

	fmt.Println("}\n\n")

	for _, section := range rootSections {
		fmt.Println(section.Typedef())
	}

	fmt.Println("\n\n")
}
