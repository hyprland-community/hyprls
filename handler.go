package hyprls

import (
	"context"
	"fmt"
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

var preferIgnoreFile = true

func (h *Handler) hasIgnoreFile() bool {
	f := ignoreFile
	stat, err := os.Stat(f)
	if err != nil {
		return false // assume failure means no file
	}
	if stat.IsDir() {
		return false
	}

	return true
}

func (h *Handler) parseIgnoresFile(uri protocol.URI) ([]string, error) {
	if !preferIgnoreFile {
		return []string{}, fmt.Errorf("preferIgnoreFile is false")
	}
	f := uri.Filename()
	if h.hasIgnoreFile() {
		b, err := os.ReadFile(f)
		if err != nil {
			h.Logger.Error("Failed to read ignore file", zap.String("file", f), zap.Error(err))
			return []string{}, err
		}
		contd := string(b)
		ignores := strings.Split(contd, "\n")
		// drop all comments and empty lines
		var filteredIgnores = make([]string, 0, len(ignores))
		for _, ig := range ignores {
			ig = strings.TrimSpace(ig)
			if ig == "" || strings.HasPrefix(ig, "#") || strings.HasPrefix(ig, "//") {
				continue
			}
			filteredIgnores = append(filteredIgnores, ig)
		}
		ignores = filteredIgnores
		h.Logger.Info("Loaded ignores from file", zap.String("file", f), zap.Strings("ignores", ignores))
		return ignores, nil
	} else {
		h.Logger.Info("Ignore file does not exist", zap.String("file", f))
		return []string{}, fmt.Errorf("ignore file does not exist")
	}
}

func (h Handler) Initialize(ctx context.Context, params *protocol.InitializeParams) (*protocol.InitializeResult, error) {
	logger = h.Logger
	workspaceFolders := params.WorkspaceFolders
	logger.Info("Loading", zap.Any("workspace", workspaceFolders))

	igs := defaultIgnores

	for _, workspaceFolder := range workspaceFolders {
		f := protocol.URI(workspaceFolder.URI + "/" + ignoreFile)
		loadedIgs, err := h.parseIgnoresFile(f)
		if err == nil {
			igs = loadedIgs
			break
		}
	}
	Ignores = igs
	logger.Info("Ignoring files", zap.Any("ignores", Ignores))

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
