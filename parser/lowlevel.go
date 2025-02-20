package parser

import (
	"errors"
	"fmt"
	"image/color"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	parser_data "github.com/hyprland-community/hyprls/parser/data"
	"go.lsp.dev/protocol"
)

var RootSection = "General"

type Position struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

func (p Position) LSP() protocol.Position {
	return protocol.Position{
		Line:      uint32(p.Line),
		Character: uint32(p.Column),
	}
}

type Assignment struct {
	Key      string   `json:"k"`
	Value    Value    `json:"v"`
	ValueRaw string   `json:"r"`
	Position Position `json:"pos"`
}

type Section struct {
	Name        string           `json:"n"`
	Start       Position         `json:"start"`
	End         Position         `json:"end"`
	Assignments []Assignment     `json:"a"`
	Variables   []CustomVariable `json:"vars"`
	Statements  []Statement      `json:"stmt"`
	Subsections []Section        `json:"sec"`
}

func (r Section) LSPRange() protocol.Range {
	return protocol.Range{
		Start: r.Start.LSP(),
		End:   r.End.LSP(),
	}
}

type CustomVariable struct {
	Assignment
}

type Statement struct {
	Keyword   Keyword  `json:"k"`
	Arguments []Value  `json:"args"`
	Position  Position `json:"pos"`
}

type Keyword string

// ValueKind represents the kind of values assignments can have.
// Reference: https://wiki.hyprland.org/Configuring/Variables/
type ValueKind int

const (
	Integer ValueKind = iota
	Bool
	Float
	Color
	Vec2
	Modmask
	String
	Gradient
	// Value contains a custom variable, its type can't be determined at compile time
	Custom
)

type ModKey int

const (
	ModShift ModKey = iota
	ModCaps
	ModControl
	ModAlt
	Mod2
	Mod3
	ModSuper
	Mod5
)

// var ColorValuePattern = regexp.MustCompile(`
//
//	(?:rgb\(([0-9a-fA-F]{6})\))|
//	(?:rgba\(([0-9a-fA-F]{8})\))|
//	(?:0x([0-9a-fA-F]{8}))
//
// `)
var ColorValuePattern = regexp.MustCompile(regexp.MustCompile(`\s+`).ReplaceAllString(`
	(?:rgb\(
		(?P<rgb_r>[0-9a-fA-F]{2})
		(?P<rgb_g>[0-9a-fA-F]{2})
		(?P<rgb_b>[0-9a-fA-F]{2})
	\))
	|
	(?:rgba\(
		(?P<rgba_r>[0-9a-fA-F]{2})
		(?P<rgba_g>[0-9a-fA-F]{2})
		(?P<rgba_b>[0-9a-fA-F]{2})
		(?P<rgba_a>[0-9a-fA-F]{2})
	\))
	|
	(?:0x
		(?P<legacy_a>[0-9a-fA-F]{2})
		(?P<legacy_r>[0-9a-fA-F]{2})
		(?P<legacy_g>[0-9a-fA-F]{2})
		(?P<legacy_b>[0-9a-fA-F]{2})
	)
`, ""))
var GradientAnglePattern = regexp.MustCompile(`(\d+)deg`)
var ModMaskSeparator = regexp.MustCompile(`[^a-zA-Z0-9,]`)

type GradientValue struct {
	Stops []Value `json:"stops,omitempty"`
	Angle float32 `json:"angle,omitempty"`
}

type Value struct {
	Kind     ValueKind     `json:"kind"`
	Bool     bool          `json:"bool"`
	Integer  int           `json:"int"`
	Float    float32       `json:"float,omitempty"`
	Color    color.RGBA    `json:"color,omitempty"`
	Vec2     [2]float32    `json:"vec2,omitempty"`
	Modmask  []ModKey      `json:"MOD,omitempty"`
	String   string        `json:"str,omitempty"`
	Gradient GradientValue `json:"gradient,omitempty"`
	Custom   string        `json:"custom,omitempty"`
	Start    Position      `json:"start"`
	End      Position      `json:"end"`
}

func (v Value) GoValue() any {
	switch v.Kind {
	case Integer:
		return v.Integer
	case Bool:
		return v.Bool
	case Float:
		return v.Float
	case Color:
		return v.Color
	case Vec2:
		return v.Vec2
	case Modmask:
		return v.Modmask
	case String:
		return v.String
	case Gradient:
		return v.Gradient
	case Custom:
		return v.Custom
	default:
		return nil
	}

}

func (k ValueKind) LSPSymbol() protocol.SymbolKind {
	switch k {
	case Integer:
		return protocol.SymbolKindNumber
	case Bool:
		return protocol.SymbolKindBoolean
	case Float:
		return protocol.SymbolKindNumber
	case Vec2:
		return protocol.SymbolKindArray
	case Modmask:
		return protocol.SymbolKindArray
	case String:
		return protocol.SymbolKindString
	default:
		return protocol.SymbolKindField
	}
}

func (v *Value) LSPRange() protocol.Range {
	return protocol.Range{
		Start: v.Start.LSP(),
		End:   v.End.LSP(),
	}
}

func (v *Value) LSPColor() protocol.Color {
	if v.Kind != Color {
		return protocol.Color{}
	}

	return protocol.Color{
		Red:   float64(v.Color.R) / 255,
		Green: float64(v.Color.G) / 255,
		Blue:  float64(v.Color.B) / 255,
		Alpha: float64(v.Color.A) / 255,
	}
}

func Parse(input string) (Section, error) {
	document := Section{
		Name:        RootSection,
		Assignments: []Assignment{},
		Subsections: []Section{},
		Start:       Position{0, 0},
	}

	sectionsStack := []*Section{&document}
	sectionDepth := 0
	endLine := 0
	for i, originalLine := range strings.Split(input, "\n") {
		currentSection := sectionsStack[sectionDepth]
		line := strings.TrimSpace(originalLine)
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		if strings.HasSuffix(line, "{") {
			sectionDepth++
			section := parseSectionStart(line)
			section.Start = Position{i, strings.Index(line, "{")}
			sectionsStack = append(sectionsStack, &section)
		}

		if strings.Contains(line, "=") {
			ass, stmt, customVar, isStatement, isCustomVar := ParseEqualLine(line, originalLine, Position{i, 0})
			pos := Position{i, strings.IndexFunc(line, unicode.IsPrint)}
			if isCustomVar {
				customVar.Position = pos
				currentSection.Variables = append(currentSection.Variables, customVar)
			} else if isStatement {
				stmt.Position = pos
				currentSection.Statements = append(currentSection.Statements, stmt)
			} else {
				ass.Position = pos
				currentSection.Assignments = append(currentSection.Assignments, ass)
			}
		}

		if line == "}" {
			currentSection.End = Position{i, strings.Index(originalLine, "}")}
			sectionsStack[sectionDepth-1].Subsections = append(sectionsStack[sectionDepth-1].Subsections, *sectionsStack[sectionDepth])
			sectionsStack = sectionsStack[:sectionDepth]
			sectionDepth--
			if sectionDepth < 0 {
				panic("unbalanced section")
			}
		}
		endLine = i
	}
	// FIXME 0 is incorrect, but do we care?
	document.End = Position{endLine, 0}
	return document, nil
}

func ParseEqualLine(line string, originalLine string, start Position) (ass Assignment, stmt Statement, customVar CustomVariable, isStatement bool, isCustomVar bool) {
	parts := strings.Split(line, "=")
	// parts[1] = strings.SplitN(parts[1], " #", 2)[0]
	// valueRaw := strings.TrimSpace(parts[1])
	key := strings.TrimSpace(parts[0])
	isCustomVar = strings.HasPrefix(key, "$")
	_, isStatement = parser_data.FindKeyword(key)

	var valueRaw string
	encounteredEquals := false
	encounteredValue := false
	valueStart := start
	valueEnd := Position{start.Line, strings.LastIndexFunc(originalLine, not(unicode.IsSpace))}
	for i, char := range originalLine {
		if !encounteredEquals && unicode.IsSpace(char) {
			continue
		}

		if char == '=' {
			encounteredEquals = true
			continue
		}

		if encounteredEquals && !encounteredValue && !unicode.IsSpace(char) {
			encounteredValue = true
			valueStart.Column = i
		}

		if encounteredValue {
			if char == '#' {
				valueEnd.Column = strings.LastIndexFunc(originalLine[:i], not(unicode.IsSpace))
				break
			}
			valueRaw += string(char)
		}
	}

	if isCustomVar {
		_ass := parseAssignment(strings.TrimPrefix(key, "$"), valueRaw, valueStart)
		_ass.Value.Start = valueStart
		_ass.Value.End = valueEnd
		customVar = CustomVariable{
			Assignment: _ass,
		}
	} else if isStatement {
		stmt = parseStatement(key, valueRaw)
	} else {
		ass = parseAssignment(key, valueRaw, valueStart)
		ass.Value.Start = valueStart
		ass.Value.End = valueEnd
	}

	return
}

func parseStatement(key string, valueRaw string) Statement {
	args := make([]Value, 0)
	for _, arg := range strings.Split(valueRaw, ",") {
		args = append(args, parseValue(strings.TrimSpace(arg) /* TODO: pass down valueStart position */, Position{0, 0}))
	}
	return Statement{
		Keyword:   Keyword(key),
		Arguments: args,
	}
}

func parseAssignment(key string, valueRaw string, valueStart Position) Assignment {
	return Assignment{
		Key:      key,
		ValueRaw: valueRaw,
		Value:    parseValue(valueRaw, valueStart),
	}
}

func parseSectionStart(line string) Section {
	return Section{
		Name:        strings.TrimSpace(strings.TrimSuffix(line, "{")),
		Subsections: []Section{},
		Assignments: []Assignment{},
	}
}

func parseValue(raw string, valueStart Position) Value {
	if strings.Contains(raw, "$") {
		return Value{
			Kind:   Custom,
			Custom: raw,
		}
	}
	if boolean, err := parseBool(raw); err == nil {
		return Value{
			Kind: Bool,
			Bool: boolean,
		}
	}

	if modmask, err := parseModMask(raw); err == nil {
		return Value{
			Kind:    Modmask,
			Modmask: modmask,
		}
	}

	if gradient, err := parseGradient(raw, valueStart); err == nil {
		return Value{
			Kind:     Gradient,
			Gradient: gradient,
		}
	}

	if color, err := ParseColor(raw); err == nil {
		return Value{
			Kind:  Color,
			Color: color,
		}
	}

	if integer, err := strconv.Atoi(raw); err == nil {
		return Value{
			Kind:    Integer,
			Integer: integer,
		}
	}

	if floating, err := strconv.ParseFloat(raw, 32); err == nil {
		return Value{
			Kind:  Float,
			Float: float32(floating),
		}
	}

	if vec, err := parseVec2(raw); err == nil {
		return Value{
			Kind: Vec2,
			Vec2: vec,
		}
	}

	return Value{
		Kind:   String,
		String: raw,
	}

}

var ModKeyNames = map[string]ModKey{
	"SHIFT":   ModShift,
	"CAPS":    ModCaps,
	"CONTROL": ModControl,
	"CTRL":    ModControl,
	"ALT":     ModAlt,
	"MOD2":    Mod2,
	"MOD3":    Mod3,
	"SUPER":   ModSuper,
	"WIN":     ModSuper,
	"LOGO":    ModSuper,
	"MOD4":    ModSuper,
	"MOD5":    Mod5,
}

func parseModMask(raw string) ([]ModKey, error) {
	potentialMods := ModMaskSeparator.Split(raw, -1)
	mods := make([]ModKey, 0, len(potentialMods))
	for _, mod := range potentialMods {
		if key, ok := ModKeyNames[strings.ToUpper(mod)]; ok {
			mods = append(mods, key)
		} else {
			return nil, fmt.Errorf("invalid mod key: %s", mod)
		}
	}
	return mods, nil
}

func parseVec2(raw string) ([2]float32, error) {
	args := strings.Split(raw, " ")
	if len(args) != 2 {
		return [2]float32{}, errors.New("invalid vec2 value")
	}

	vec := [2]float32{}
	for i, arg := range args {
		f, err := strconv.ParseFloat(arg, 32)
		if err != nil {
			return [2]float32{}, errors.New("invalid vec2 value")
		}
		vec[i] = float32(f)
	}

	return vec, nil
}

func ParseColor(raw string) (color.RGBA, error) {
	if ColorValuePattern.MatchString(raw) {
		matches := ColorValuePattern.FindStringSubmatch(raw)
		return color.RGBA{
			R: decodeHexComponent(matches, "r", 0),
			G: decodeHexComponent(matches, "g", 0),
			B: decodeHexComponent(matches, "b", 0),
			A: decodeHexComponent(matches, "a", 0xff),
		}, nil
	}
	return color.RGBA{0, 0, 0, 0}, errors.New("invalid color value")
}

func parseGradient(raw string, valueStart Position) (GradientValue, error) {
	args := strings.Split(raw, " ")
	value := GradientValue{}

	cursorAt := 0
	for i, arg := range args {
		originalArg := arg
		arg = strings.TrimSpace(arg)
		// +1 since we splitted the spaces out
		cursorAt += len(originalArg) - len(arg) + 1

		color, err := ParseColor(arg)
		if err != nil {
			if i == len(args)-1 {
				if !GradientAnglePattern.MatchString(arg) {
					return GradientValue{}, errors.New("invalid gradient angle")
				}
				angle, _ := strconv.ParseFloat(GradientAnglePattern.FindStringSubmatch(arg)[1], 32)
				value.Angle = float32(angle)
				return value, nil
			} else {
				return GradientValue{}, errors.New("invalid gradient value")
			}
		}

		value.Stops = append(value.Stops, Value{
			Kind:  Color,
			Color: color,
			Start: Position{valueStart.Line, valueStart.Column + cursorAt - 1},
			End:   Position{valueStart.Line, valueStart.Column + cursorAt + len(arg) - 1},
		})
		cursorAt += len(arg)
	}

	return value, nil
}

func hexToColor(hexstring string) color.RGBA {
	components := []uint64{0, 0, 0, 0xff}
	for i := 0; i < 4; i++ {
		if len(hexstring) <= i*2 {
			components[i], _ = strconv.ParseUint(hexstring[i*2:i*2+2], 16, 8)
		}
	}
	return color.RGBA{uint8(components[0]), uint8(components[1]), uint8(components[2]), uint8(components[3])}
}

func decodeHexComponent(matches []string, component string, fallback uint8) uint8 {
	toDecode := matches[ColorValuePattern.SubexpIndex("rgba_"+component)]
	if component != "a" && toDecode == "" {
		toDecode = matches[ColorValuePattern.SubexpIndex("rgb_"+component)]
	}
	if toDecode == "" {
		toDecode = matches[ColorValuePattern.SubexpIndex("legacy_"+component)]
	}

	if toDecode == "" {
		return fallback
	}

	decoded, _ := strconv.ParseUint(toDecode, 16, 8)
	return uint8(decoded)
}

func parseBool(raw string) (bool, error) {
	switch strings.TrimSpace(raw) {
	case "true", "yes", "on", "1":
		return true, nil
	case "false", "no", "off", "0":
		return false, nil
	default:
		return false, errors.New("invalid bool value")
	}
}

func (s Section) WalkValues(f func(assignment *Assignment, v *Value)) {
	for _, a := range s.Assignments {
		f(&a, &a.Value)
	}
	for _, s := range s.Subsections {
		s.WalkValues(func(a *Assignment, v *Value) {
			f(a, v)
		})
	}
}

func (s Section) WalkCustomVariables(f func(v *CustomVariable)) {
	for _, v := range s.Variables {
		f(&v)
	}
	for _, s := range s.Subsections {
		s.WalkCustomVariables(f)
	}
}

func ValueKindFromString(s string) (ValueKind, error) {
	switch strings.ToLower(s) {
	case "int", "integer":
		return Integer, nil
	case "bool":
		return Bool, nil
	case "float":
		return Float, nil
	case "color":
		return Color, nil
	case "vec2":
		return Vec2, nil
	case "mod", "modmask":
		return Modmask, nil
	case "str", "string":
		return String, nil
	case "gradient":
		return Gradient, nil
	default:
		return Custom, fmt.Errorf("unknown value kind: %s", s)
	}

}
