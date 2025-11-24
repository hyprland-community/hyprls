package hyprls

import (
	"context"
	"strings"

	"go.lsp.dev/protocol"
	"go.uber.org/zap"
)

func extractIgnoresFromChangeConfigSettings(params *protocol.DidChangeConfigurationParams) (ignores []string, isAvailable bool) {
	settings, ok := (params.Settings).(map[string]any)
	if !ok {
		return nil, false
	}
	hyprls, ok := settings["hyprls"].(map[string]any)
	if !ok {
		return nil, false
	}
	if pfIgnore, ok := hyprls["preferIgnoreFile"]; ok {
		if pfIgnoreBool, ok := pfIgnore.(bool); ok {
			preferIgnoreFile = pfIgnoreBool
		}
	}
	if ignore, ok := hyprls["ignore"]; ok {
		if arr, ok := ignore.([]any); ok {
			isAvailable = true
			for _, v := range arr {
				if s, ok := v.(string); ok {
					ignores = append(ignores, s)
				}
			}
		}
	}
	return
}

func (h Handler) DidChangeConfiguration(ctx context.Context, params *protocol.DidChangeConfigurationParams) error {
	newIgnores, updated := extractIgnoresFromChangeConfigSettings(params)
	if updated && !preferIgnoreFile {
		Ignores = newIgnores
		h.Logger.Info("configuration changed", zap.Strings("ignores", Ignores))
	}
	return nil
}

func (h Handler) DidChangeWatchedFiles(ctx context.Context, params *protocol.DidChangeWatchedFilesParams) error {
	for _, change := range params.Changes {
		if strings.HasSuffix(string(change.URI), ignoreFile) {
			if !preferIgnoreFile {
				h.Logger.Info("Ignoring change to ignore file as preferIgnoreFile is false")
				continue
			}
			h.Logger.Info("Ignore file changed", zap.String("uri", string(change.URI)))
			if change.Type == protocol.FileChangeTypeDeleted {
				Ignores = []string{"hyprlock.conf", "hypridle.conf"}
				h.Logger.Info("Ignore file deleted, resetting to defaults")
				continue
			}
			newIgnores, err := h.parseIgnoresFile(change.URI)
			if err == nil {
				Ignores = newIgnores
				h.Logger.Info("Updated ignores from changed ignore file", zap.Strings("ignores", Ignores))
			}
		}
	}
	return nil
}
