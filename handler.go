package hyprls

import (
	"context"

	"go.lsp.dev/protocol"
	"go.uber.org/zap"
)

type Handler struct {
	Server protocol.Server
	Logger *zap.Logger
}

func NewHandler(ctx context.Context, server protocol.Server, logger *zap.Logger) (Handler, context.Context, error) {

	return Handler{
		Server: server,
		Logger: logger,
	}, context.WithValue(ctx, "state", state{}), nil
}

func (h Handler) Initialize(ctx context.Context, params *protocol.InitializeParams) (*protocol.InitializeResult, error) {
	logger = h.Logger
	return &protocol.InitializeResult{
		Capabilities: protocol.ServerCapabilities{
			HoverProvider:          true,
			DocumentSymbolProvider: true,
			ColorProvider:          true,
			CompletionProvider: &protocol.CompletionOptions{
				ResolveProvider:   false,
				TriggerCharacters: []string{},
			},
			TextDocumentSync: protocol.TextDocumentSyncOptions{
				OpenClose: true,
				Change:    protocol.TextDocumentSyncKindFull,
			},
			Workspace: &protocol.ServerCapabilitiesWorkspace{
				WorkspaceFolders: &protocol.ServerCapabilitiesWorkspaceFolders{
					Supported:           true,
					ChangeNotifications: true,
				},
			},
		},
		ServerInfo: &protocol.ServerInfo{
			Name:    "hyprls",
			Version: Version,
		},
	}, nil
}

func (h Handler) Initialized(ctx context.Context, params *protocol.InitializedParams) error {
	return nil
}

func (h Handler) Shutdown(ctx context.Context) error {
	return nil
}

func (h Handler) Exit(ctx context.Context) error {
	return nil
}
