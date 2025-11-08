package hyprls

import (
	"context"
	"os"
	"strings"

	"go.lsp.dev/protocol"
	"go.uber.org/zap"
)

type Handler struct {
	Server protocol.Server
	Logger *zap.Logger
}

type GlobalContextKey string

func NewHandler(ctx context.Context, server protocol.Server, logger *zap.Logger) (Handler, context.Context, error) {

	return Handler{
		Server: server,
		Logger: logger,
	}, context.WithValue(ctx, GlobalContextKey("state"), state{}), nil
}

const ignoreFile = ".hyprlsignore"

func (h Handler) Initialize(ctx context.Context, params *protocol.InitializeParams) (*protocol.InitializeResult, error) {
	logger = h.Logger
	workspaceFolders := params.WorkspaceFolders
	logger.Info("Loading", zap.Any("workspace", workspaceFolders))

	for _, workspaceFolder := range workspaceFolders {
		f := protocol.URI(workspaceFolder.URI + "/" + ignoreFile).Filename()
		if info, err := os.Stat(f); err == nil && !info.IsDir() {
			logger.Debug("Loading ignore file" + ignoreFile)
			b, err := os.ReadFile(f)
			if err == nil {
				contd := string(b)
				ignores = strings.Split(contd, "\n")
			}
		}
	}
	logger.Debug("Ignoring files", zap.Any("ignores", ignores))

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
