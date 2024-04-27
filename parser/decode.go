package parser

import (
	"fmt"
	"reflect"

	parser_data "github.com/ewen-lbh/hyprls/parser/data"
)

func (root Section) Decode() (Configuration, error) {
	config := Configuration{
		CustomVariables: make(map[string]string, 0),
	}
	root.WalkCustomVariables(func(v *CustomVariable) {
		config.CustomVariables[v.Key] = v.ValueRaw
	})

	for _, ass := range root.Assignments {
		def := parser_data.FindVariableDefinitionInSection("General", ass.Key)
		if def == nil {
			availableKeys := make([]string, 0)
			for _, v := range parser_data.FindSectionDefinitionByName("General").Variables{
				availableKeys = append(availableKeys, v.Name)
			}
			return Configuration{}, fmt.Errorf("unknown variable General > %s. Available keys are %v", ass.Key, availableKeys)
		}
		fmt.Printf("adding %s=%#v to .General", def.PascalCaseName(), ass.Value.GoValue())
		setValue(&config, def.PascalCaseName(), ass.Value.GoValue())
	}

	return config, nil
}

func setValue(obj any, field string, value any) {
	ref := reflect.ValueOf(obj)

	// if its a pointer, resolve its value
	if ref.Kind() == reflect.Ptr {
		ref = reflect.Indirect(ref)
	}

	if ref.Kind() == reflect.Interface {
		ref = ref.Elem()
	}

	// should double check we now have a struct (could still be anything)
	if ref.Kind() != reflect.Struct {
		panic("cannot setValue on a non-struct")
	}

	prop := ref.FieldByName(field)
	prop.Set(reflect.ValueOf(value))
}
