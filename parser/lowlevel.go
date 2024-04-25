package parser

import (
	"errors"
	"fmt"
	"image/color"
	"regexp"
	"strconv"
	"strings"
)

var RootSection = "____root____"

type Assignment struct {
	Key      string `json:"k"`
	Value    Value  `json:"v"`
	ValueRaw string `json:"r"`
	Line     int    `json:"l"`
}

type Section struct {
	Name        string           `json:"n"`
	StartLine   int              `json:"sl"`
	EndLine     int              `json:"el"`
	Assignments []Assignment     `json:"a"`
	Variables   []CustomVariable `json:"vars"`
	Statements  []Statement      `json:"stmt"`
	Subsections []Section        `json:"sec"`
}

type CustomVariable struct {
	Assignment
}

type Statement struct {
	Keyword   Keyword `json:"k"`
	Arguments []Value `json:"args"`
	Line      int     `json:"l"`
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

var ColorValuePattern = regexp.MustCompile(`
	(?:rgb\(([0-9a-fA-F]{6})\))|
	(?:rgba\(([0-9a-fA-F]{8})\))|
	(?:0x([0-9a-fA-F]{8}))
`)
var GradientAnglePattern = regexp.MustCompile(`(\d+)deg`)
var ModMaskSeparator = regexp.MustCompile(`[^a-zA-Z0-9,]`)

type GradientValue struct {
	Stops []color.RGBA `json:"stops,omitempty"`
	Angle float32      `json:"angle,omitempty"`
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
			ass, stmt, customVar, isStatement, isCustomVar := parseEqualLine(line)
			if isCustomVar {
				customVar.Line = i + 1
				currentSection.Variables = append(currentSection.Variables, customVar)
			} else if isStatement {
				stmt.Line = i + 1
				currentSection.Statements = append(currentSection.Statements, stmt)
			} else {
				ass.Line = i + 1
				currentSection.Assignments = append(currentSection.Assignments, ass)
			}
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

func parseEqualLine(line string) (ass Assignment, stmt Statement, customVar CustomVariable, isStatement bool, isCustomVar bool) {
	parts := strings.Split(line, "=")
	parts[1] = strings.SplitN(parts[1], " #", 2)[0]
	valueRaw := strings.TrimSpace(parts[1])
	key := strings.TrimSpace(parts[0])
	isCustomVar = strings.HasPrefix(key, "$")
	isStatement = IsKeyword(key)

	if isCustomVar {
		_ass := parseAssignment(strings.TrimPrefix(key, "$"), valueRaw)
		customVar = CustomVariable{
			Assignment: _ass,
		}
	} else if isStatement {
		stmt = parseStatement(key, valueRaw)
	} else {
		ass = parseAssignment(key, valueRaw)
	}

	return
}

func parseStatement(key string, valueRaw string) Statement {
	args := make([]Value, 0)
	for _, arg := range strings.Split(valueRaw, ",") {
		args = append(args, parseValue(strings.TrimSpace(arg)))
	}
	return Statement{
		Keyword:   Keyword(key),
		Arguments: args,
	}
}

func parseAssignment(key string, valueRaw string) Assignment {
	return Assignment{
		Key:      key,
		ValueRaw: valueRaw,
		Value:    parseValue(valueRaw),
	}
}

func parseSectionStart(line string) Section {
	return Section{
		Name:        strings.TrimSpace(strings.TrimSuffix(line, "{")),
		Subsections: []Section{},
		Assignments: []Assignment{},
	}
}

func parseValue(raw string) Value {
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

	if gradient, err := parseGradient(raw); err == nil {
		return Value{
			Kind:     Gradient,
			Gradient: gradient,
		}
	}

	if color, err := parseColor(raw); err == nil {
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

func parseModMask(raw string) ([]ModKey, error) {
	potentialMods := ModMaskSeparator.Split(raw, -1)
	mods := make([]ModKey, 0, len(potentialMods))
	for _, mod := range potentialMods {
		switch mod {
		case "SHIFT":
			mods = append(mods, ModShift)
		case "CAPS":
			mods = append(mods, ModCaps)
		case "CONTROL", "CTRL":
			mods = append(mods, ModControl)
		case "ALT":
			mods = append(mods, ModAlt)
		case "MOD2":
			mods = append(mods, Mod2)
		case "MOD3":
			mods = append(mods, Mod3)
		case "SUPER", "WIN", "LOGO", "MOD4":
			mods = append(mods, ModSuper)
		case "MOD5":
			mods = append(mods, Mod5)
		default:
			return mods, fmt.Errorf("invalid mod %s", mod)
		}
	}
	return mods, nil
}

func parseVec2(raw string) ([2]float32, error) {
	args := strings.Split(raw, ",")
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

func parseColor(raw string) (color.RGBA, error) {
	if ColorValuePattern.MatchString(raw) {
		matches := ColorValuePattern.FindStringSubmatch(raw)
		return hexToColor(matches[1]), nil
	}
	return color.RGBA{0, 0, 0, 0}, errors.New("invalid color value")
}

func parseGradient(raw string) (GradientValue, error) {
	args := strings.Split(raw, " ")
	value := GradientValue{}

	for i, arg := range args {
		arg = strings.TrimSpace(arg)
		color, err := parseColor(arg)
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

		value.Stops = append(value.Stops, color)
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

func parseBool(raw string) (bool, error) {
	switch raw {
	case "true", "yes", "on", "1":
		return true, nil
	case "false", "no", "off", "0":
		return false, nil
	default:
		return false, errors.New("invalid bool value")
	}
}
