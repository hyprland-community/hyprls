package main

import (
	"strings"
)

var RootSection = "____root____"

type Assignment struct {
	Key   string `json:"k"`
	Value string `json:"v"`
	Line  int    `json:"l"`
}

type Section struct {
	Name        string       `json:"n"`
	StartLine   int          `json:"sl"`
	EndLine     int          `json:"el"`
	Assignments []Assignment `json:"a"`
	Subsections []Section    `json:"s"`
}

func Parse(input string) (Section, error) {
	document := Section{
		Name:        RootSection,
		Assignments: []Assignment{},
		Subsections: []Section{},
		StartLine:   1,
	}

	sectionsStack := []*Section{&document}
	sectionDepth := 0
	endLine := 0
	for i, line := range strings.Split(input, "\n") {
		currentSection := sectionsStack[sectionDepth]
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		if strings.HasSuffix(line, "{") {
			sectionDepth++
			section := parseSectionStart(line)
			section.StartLine = i + 1
			sectionsStack = append(sectionsStack, &section)
		}

		if strings.Contains(line, "=") {
			ass := parseAssignment(line)
			ass.Line = i + 1
			currentSection.Assignments = append(currentSection.Assignments, ass)
		}

		if line == "}" {
			currentSection.EndLine = i + 1
			sectionsStack[sectionDepth-1].Subsections = append(sectionsStack[sectionDepth-1].Subsections, *sectionsStack[sectionDepth])
			sectionsStack = sectionsStack[:sectionDepth]
			sectionDepth--
			if sectionDepth < 0 {
				panic("unbalanced section")
			}
		}
		endLine = i
	}
	document.EndLine = endLine + 1
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
