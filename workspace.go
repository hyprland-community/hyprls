package hyprls

import (
	"context"

	"go.lsp.dev/protocol"
)

func extractExcludes(params *protocol.DidChangeConfigurationParams) (excludes []string) {
	if settings, ok := (params.Settings).(map[string]any); ok {
		if exclude, ok := settings["exclude"]; ok {
			if s, ok := exclude.(string); ok {
				excludes = append(excludes, s)
			}
			if arr, ok := exclude.([]any); ok {
				for _, e := range arr {
					if s, ok := e.(string); ok {
						excludes = append(excludes, s)
					}
				}
			}
		}
	}
	return
}

func (h Handler) DidChangeConfiguration(ctx context.Context, params *protocol.DidChangeConfigurationParams) error {
	ignores = extractExcludes(params)
	return nil
}
