package parser

import "strings"

var knownKeywords = []string{
	"env", "envd", "monitor", "bind", "unbind", "submap", "windowrule", "windowrulev2", "layerrule", "workspace", "animation", "bezier", "exec", "exec-once",
}

var keywordsWithFlags = map[string][]string{
	"bind": {"r", "l", "e", "n", "m", "t", "i"},
}

func IsKeyword(key string) bool {
	for _, k := range knownKeywords {
		if k == key {
			return true
		}

		if strings.HasPrefix(key, k) {
			if availableFlags, ok := keywordsWithFlags[k]; ok {
				usedFlags := strings.Split(strings.TrimPrefix(key, k), "")
				for _, used := range usedFlags {
					flagIsValid := false
					for _, available := range availableFlags {
						if used == available {
							flagIsValid = true
							break
						}
					}
					if !flagIsValid {
						return false
					}
				}
				return true
			}
		}
	}
	return false
}
