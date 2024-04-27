package parser_data

import (
	"fmt"
	"strings"
)

func FindSectionDefinitionByName(name string) *SectionDefinition {
	for _, sec := range Sections {
		if sec.Name() == name || sec.JSONName() == name {
			return &sec
		}
	}
	return nil
}

type SectionDefinition struct {
	Path        []string
	Subsections []SectionDefinition
	Variables   []VariableDefinition
}

func (s SectionDefinition) Name() string {
	if len(s.Path) == 0 {
		return ""
	}
	return s.Path[len(s.Path)-1]
}

func (s SectionDefinition) JSONName() string {
	return strings.ToLower(s.Name())
}

func (s SectionDefinition) TypeName() string {
	return "Configuration" + toPascalCase(strings.Join(s.Path, "_"))
}

func (s SectionDefinition) Typedef() string {
	out := fmt.Sprintf("type %s struct {\n", s.TypeName())
	for _, def := range s.Variables {
		out += fmt.Sprintf("\t// %s\n", def.Description)
		out += fmt.Sprintf("\t%s %s `json:\"%s\"`\n", def.PascalCaseName(), def.GoType(), def.Name)
		out += "\n"
	}
	for _, sec := range s.Subsections {
		out += fmt.Sprintf("\t%s %s `json:\"%s\"`\n", sec.Name(), sec.TypeName(), strings.ToLower(sec.Name()))
	}
	out += "}\n"
	out += "\n"

	for _, sec := range s.Subsections {
		out += sec.Typedef()
	}
	return out
}
