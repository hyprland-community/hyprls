package parser_data

import (
	"fmt"
	"strings"
)

type KeywordDefinition struct {
	Name                     string
	Description              string
	documentationHeadingSlug string
	documentationFile        string
	Flags                    []string
}

func (k KeywordDefinition) DocumentationLink() string {
	return fmt.Sprintf("https://wiki.hyprland.org/Configuring/%s/#%s", k.documentationFile, k.documentationHeadingSlug)
}

var Keywords = []KeywordDefinition{
	{
		Name:                     "submap",
		documentationHeadingSlug: "submaps",
		documentationFile:        "Binds",
		Flags:                    []string{},
	},
	{
		Name:                     "windowrule",
		documentationHeadingSlug: "window-rules",
		documentationFile:        "Window-Rules",
		Flags:                    []string{},
	},
	{
		Name:                     "windowrulev2",
		documentationHeadingSlug: "syntax",
		documentationFile:        "Window-Rules",
		Flags:                    []string{},
	},
	{
		Name:                     "layerrule",
		documentationHeadingSlug: "layer-rules",
		documentationFile:        "Window-Rules",
		Flags:                    []string{},
	},
	{
		Name:                     "workspace",
		documentationHeadingSlug: "rules",
		documentationFile:        "Workspace-Rules",
		Flags:                    []string{},
	},
	{
		Name:                     "animation",
		documentationHeadingSlug: "general",
		documentationFile:        "Animations",
		Flags:                    []string{},
	},
	{
		Name:                     "bezier",
		documentationHeadingSlug: "curves",
		documentationFile:        "Animations",
		Flags:                    []string{},
	},
	{
		Name:                     "exec",
		documentationHeadingSlug: "executing",
		documentationFile:        "Keywords",
		Flags:                    []string{},
	},
	{
		Name:                     "exec-once",
		documentationHeadingSlug: "executing",
		documentationFile:        "Keywords",
		Flags:                    []string{},
	},
	{
		Name:                     "source",
		documentationHeadingSlug: "sourcing-multi-file",
		documentationFile:        "Keywords",
		Flags:                    []string{},
	},
	{
		Name:                     "env",
		documentationHeadingSlug: "setting-the-environment",
		documentationFile:        "Keywords",
		Flags:                    []string{"d"},
	},
	{
		Name:                     "monitor",
		documentationHeadingSlug: "general",
		documentationFile:        "Monitors",
		Flags:                    []string{},
	},
	{
		Name:                     "bind",
		documentationHeadingSlug: "basic",
		documentationFile:        "Binds",
		Flags:                    []string{"r", "l", "e", "n", "m", "t", "i"},
	},
	{
		Name:                     "unbind",
		documentationHeadingSlug: "unbind",
		documentationFile:        "Binds",
		Flags:                    []string{},
	},
}

func FindKeyword(key string) (keyword KeywordDefinition, found bool) {
	for _, k := range Keywords {
		if key == k.Name {
			return k, true
		}

		if strings.HasPrefix(key, k.Name) {
			if len(k.Flags) > 0 {
				usedFlags := strings.Split(strings.TrimPrefix(key, k.Name), "")
				for _, used := range usedFlags {
					flagIsValid := false
					for _, available := range k.Flags {
						if used == available {
							flagIsValid = true
							break
						}
					}
					if !flagIsValid {
						return KeywordDefinition{}, false
					}
				}
				return k, true
			}
		}
	}
	return KeywordDefinition{}, false
}
