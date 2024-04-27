package parser_data

import (
	"bytes"
	_ "embed"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/anaskhan96/soup"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

var md = goldmark.New(goldmark.WithExtensions(extension.GFM))

func debug(msg string, fmtArgs ...any) {
	// fmt.Fprintf(os.Stderr, msg, fmtArgs...)
}

//go:embed Variables.md
var documentationSource []byte

//go:embed Master-Layout.md
var masterLayoutDocumentationSource []byte

//go:embed Dwindle-Layout.md
var dwindleLayoutDocumentationSource []byte

var Sections = []SectionDefinition{}

func (s SectionDefinition) VariableDefinition(name string) *VariableDefinition {
	for _, v := range s.Variables {
		if v.Name == name {
			return &v
		}
	}
	return nil
}

func init() {
	Sections = parseDocumentationMarkdown(documentationSource, 3)
	Sections = append(Sections, parseDocumentationMarkdownWithRootSectionName(masterLayoutDocumentationSource, 2, "Master")...)
	Sections = append(Sections, parseDocumentationMarkdownWithRootSectionName(dwindleLayoutDocumentationSource, 2, "Dwindle")...)
}

func parseDocumentationMarkdownWithRootSectionName(source []byte, headingRootLevel int, rootSectionName string) []SectionDefinition {
	sections := parseDocumentationMarkdown(source, headingRootLevel)
	for i := range sections {
		sections[i].Path[0] = rootSectionName
	}
	return sections
}

func parseDocumentationMarkdown(source []byte, headingRootLevel int) (sections []SectionDefinition) {
	var html bytes.Buffer
	err := md.Convert(source, &html)
	if err != nil {
		panic(err)
	}

	document := soup.HTMLParse(html.String())
	for _, table := range document.FindAll("table") {
		if !arraysEqual(tableHeaderCells(table), []string{"name", "description", "type", "default"}) {
			continue
		}

		// fmt.Printf("Processing table %s\n", table.HTML())
		section := SectionDefinition{
			Path: tablePath(table, headingRootLevel),
		}
		section.Variables = make([]VariableDefinition, 0)
		for _, row := range table.FindAll("tr")[1:] {
			cells := row.FindAll("td")
			if len(cells) != 4 {
				continue
			}

			section.Variables = append(section.Variables, VariableDefinition{
				Name:        cells[0].FullText(),
				Description: cells[1].FullText(),
				Type:        cells[2].FullText(),
				Default:     cells[3].FullText()})
		}
		sections = append(sections, section)
	}

	for i, section := range sections {
		if len(section.Path) == 1 {
			sections[i] = section.AttachSubsections(sections)
		}
	}
	return sections
}

func (s SectionDefinition) AttachSubsections(sections []SectionDefinition) SectionDefinition {
	// TODO make it work for recursively nested sections
	s.Subsections = make([]SectionDefinition, 0)
	for _, section := range sections {
		if len(section.Path) == 1 {
			continue
		}
		if section.Path[0] == s.Name() {
			debug("adding %s to %s\n", section.Name(), s.Name())
			s.Subsections = append(s.Subsections, section)
		}
	}
	return s
}

func tableHeaderCells(table soup.Root) []string {
	headerCells := table.FindAll("th")
	cells := make([]string, 0, len(headerCells))
	for _, cell := range headerCells {
		cells = append(cells, cell.FullText())
	}
	return cells
}

func tablePath(table soup.Root, headingRootLevel int) []string {
	header := backtrackToNearestHeader(table)
	level, err := strconv.Atoi(header.NodeValue[1:])
	if err != nil {
		panic(err)
	}
	if level <= headingRootLevel {
		return []string{header.FullText()}
	}
	return append(tablePath(header.FindPrevElementSibling(), headingRootLevel), header.FullText())
}

func backtrackToNearestHeader(element soup.Root) soup.Root {
	if element.NodeValue != "table" {
		debug("backtracking to nearest header from %s\n", element.HTML())
	}
	if regexp.MustCompile(`^h[1-6]$`).MatchString(element.NodeValue) {
		debug("-> returning from backtrack with %s\n", element.HTML())
		return element
	}
	prev := element.FindPrevElementSibling()
	debug("-> prev is %s\n", prev.HTML())
	return backtrackToNearestHeader(prev)
}

var MarkdownHeaderPattern = regexp.MustCompile(`^(#+)\s+(.*)$`)
var MarkdownTableStart = `| name | description | type | default |`

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

type VariableDefinition struct {
	Name        string
	Description string
	Type        string
	Default     string
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

func (v VariableDefinition) PrettyDefault() string {
	if v.Default == "[[Empty]]" {
		return "*(empty)*"
	}
	return v.Default
}

func (v VariableDefinition) GoType() string {
	switch v.Type {
	case "int":
		return "int"
	case "bool":
		return "bool"
	case "float", "floatvalue":
		return "float32"
	case "color":
		return "color.RGBA"
	case "vec2":
		return "[2]float32"
	case "MOD":
		return "[]ModKey"
	case "str", "string":
		return "string"
	case "gradient":
		return "GradientValue"
	default:
		panic("unknown type: " + v.Type)
	}

}

func (v VariableDefinition) ParserType() string {
	// Integer ValueKind = iota
	// Bool
	// Float
	// Color
	// Vec2
	// Modmask
	// String
	// Gradient
	// Custom

	switch v.Type {
	case "int":
		return "Integer"
	case "bool":
		return "Bool"
	case "float":
		return "Float"
	case "color":
		return "Color"
	case "vec2":
		return "Vec2"
	case "MOD":
		return "Modmask"
	case "str", "string":
		return "String"
	case "gradient":
		return "Gradient"
	default:
		panic("unknown type: " + v.Type)
	}
}

func (v VariableDefinition) PascalCaseName() string {
	return toPascalCase(v.Name)
}

func arraysEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if strings.TrimSpace(v) != strings.TrimSpace(b[i]) {
			return false
		}
	}
	return true
}

func toPascalCase(s string) string {
	out := ""
	for _, word := range regexp.MustCompile(`[-_\.]`).Split(s, -1) {
		out += strings.ToUpper(word[:1]) + word[1:]
	}
	return out
}
