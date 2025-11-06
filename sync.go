package hyprls

import (
	"context"
	"errors"

	"go.lsp.dev/protocol"
	"go.uber.org/zap"
)

func (h Handler) DidChange(ctx context.Context, params *protocol.DidChangeTextDocumentParams) error {
	if isExclude(params.TextDocument.URI) {
		return nil
	}
	logger.Debug("LSP:DidChange", zap.Any("params", params))
	if len(params.ContentChanges) > 0 {
		openedFiles[params.TextDocument.URI] = params.ContentChanges[len(params.ContentChanges)-1].Text
	}
	return nil
}

func (h Handler) DidClose(ctx context.Context, params *protocol.DidCloseTextDocumentParams) error {
	if isExclude(params.TextDocument.URI) {
		return nil
	}
	delete(openedFiles, params.TextDocument.URI)
	return nil
}

func (h Handler) DidOpen(ctx context.Context, params *protocol.DidOpenTextDocumentParams) error {
	if isExclude(params.TextDocument.URI) {
		return nil
	}
	file(params.TextDocument.URI)
	return nil
}

func (h Handler) DidSave(ctx context.Context, params *protocol.DidSaveTextDocumentParams) error {
	return errors.New("unimplemented")
}
