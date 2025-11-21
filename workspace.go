package hyprls

import (
	"context"

	"go.lsp.dev/protocol"
	"go.uber.org/zap"
)

func extractIgnores(params *protocol.DidChangeConfigurationParams) (ignores []string) {
	if settings, ok := (params.Settings).(map[string]any); ok {
		if hyprls, ok := settings["hyprls"].(map[string]any); ok {
			settings = hyprls
		}

		if ignore, ok := settings["ignore"]; ok {
			if arr, ok := ignore.([]any); ok {
				for _, v := range arr {
					if s, ok := v.(string); ok {
						ignores = append(ignores, s)
					}
				}
			}
		}
	}
	return
}

func (h Handler) DidChangeConfiguration(ctx context.Context, params *protocol.DidChangeConfigurationParams) error {
	Ignores = extractIgnores(params)
	h.Logger.Info("configuration changed", zap.Strings("ignores", Ignores))
	return nil
}
