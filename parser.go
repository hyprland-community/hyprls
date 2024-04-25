package main

import (
	"strings"
)

var RootSection = "____root____"

type Assignment struct {
	Key   string `json:"k"`
	Value string `json:"v"`
}

type Section struct {
	Name        string       `json:"n"`
	Assignments []Assignment `json:"a"`
	Subsections []Section    `json:"s"`
}

func Parse(input string) (Section, error) {
	document := Section{
		Name:        RootSection,
		Assignments: []Assignment{},
		Subsections: []Section{},
	}

	sectionsStack := []*Section{&document}
	sectionDepth := 0
	for _, line := range strings.Split(input, "\n") {
		currentSection := sectionsStack[sectionDepth]
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		if strings.HasSuffix(line, "{") {
			sectionDepth++
			section := parseSectionStart(line)
			sectionsStack = append(sectionsStack, &section)
		}

		if strings.Contains(line, "=") {
			currentSection.Assignments = append(currentSection.Assignments, parseAssignment(line))
		}

		if line == "}" {
			sectionsStack[sectionDepth-1].Subsections = append(sectionsStack[sectionDepth-1].Subsections, *sectionsStack[sectionDepth])
			sectionsStack = sectionsStack[:sectionDepth]
			sectionDepth--
			if sectionDepth < 0 {
				panic("unbalanced section")
			}
		}

	}
	return document, nil
}

func parseAssignment(line string) Assignment {
	parts := strings.Split(line, "=")
	parts[1] = strings.SplitN(parts[1], " #", 2)[0]
	return Assignment{
		Key:   strings.TrimSpace(parts[0]),
		Value: strings.TrimSpace(parts[1]),
	}
}

func parseSectionStart(line string) Section {
	return Section{
		Name:        strings.TrimSpace(strings.TrimSuffix(line, "{")),
		Subsections: []Section{},
		Assignments: []Assignment{},
	}
}
