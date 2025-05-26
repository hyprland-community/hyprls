package parser_data

func FindVariableDefinitionInSection(sectionName, variableName string) *VariableDefinition {
	sec := FindSectionDefinitionByName(sectionName)
	if sec == nil {
		return nil
	}
	return sec.VariableDefinition(variableName)
}

type VariableDefinition struct {
	Name        string
	Description string
	Type        string
	Default     string
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
	case "font_weight":
		return "uint8"
	default:
		panic("unknown type: " + v.Type)
	}

}



func (v VariableDefinition) ParserTypeString() string {
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
