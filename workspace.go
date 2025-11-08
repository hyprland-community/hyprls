package hyprls

import (
	"context"

	"go.lsp.dev/protocol"
)

func extractIncludes(params *protocol.DidChangeConfigurationParams) (extractedIgnores []string) {
	if settings, ok := (params.Settings).(map[string]any); ok {
		if ignore, ok := settings["ignore"]; ok {
			if s, ok := ignore.(string); ok {
				extractedIgnores = append(extractedIgnores, s)
			}
			if arr, ok := ignore.([]any); ok {
				for _, e := range arr {
					if s, ok := e.(string); ok {
						extractedIgnores = append(extractedIgnores, s)
					}
				}
			}
		}
	}
	return
}

func (h Handler) DidChangeConfiguration(ctx context.Context, params *protocol.DidChangeConfigurationParams) error {
	ignores = extractIncludes(params)
	return nil
}
