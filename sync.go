package hyprls

import (
	"context"
	"errors"

	"go.lsp.dev/protocol"
	"go.uber.org/zap"
)


func (h Handler) DidChange(ctx context.Context, params *protocol.DidChangeTextDocumentParams) error {
	logger.Debug("LSP:DidChange", zap.Any("params", params))
	openedFiles[params.TextDocument.URI] = params.ContentChanges[len(params.ContentChanges)-1].Text
	return nil
}

func (h Handler) DidClose(ctx context.Context, params *protocol.DidCloseTextDocumentParams) error {
	delete(openedFiles, params.TextDocument.URI)
	return nil
}

func (h Handler) DidOpen(ctx context.Context, params *protocol.DidOpenTextDocumentParams) error {
	file(params.TextDocument.URI)
	return nil
}

func (h Handler) DidSave(ctx context.Context, params *protocol.DidSaveTextDocumentParams) error {
	return errors.New("unimplemented")
}
