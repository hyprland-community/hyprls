package hyprls

import (
	"context"
	"errors"

	"go.lsp.dev/protocol"
)

func (h Handler) Definition(ctx context.Context, params *protocol.DefinitionParams) ([]protocol.Location, error) {
	return []protocol.Location{}, errors.New("unimplemented")
}

func (h Handler) WorkDoneProgressCancel(ctx context.Context, params *protocol.WorkDoneProgressCancelParams) error {
	return errors.New("unimplemented")
}

func (h Handler) LogTrace(ctx context.Context, params *protocol.LogTraceParams) error {
	return errors.New("unimplemented")
}

func (h Handler) SetTrace(ctx context.Context, params *protocol.SetTraceParams) error {
	return errors.New("unimplemented")
}

func (h Handler) CodeAction(ctx context.Context, params *protocol.CodeActionParams) ([]protocol.CodeAction, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) CodeLens(ctx context.Context, params *protocol.CodeLensParams) ([]protocol.CodeLens, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) CodeLensResolve(ctx context.Context, params *protocol.CodeLens) (*protocol.CodeLens, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) Declaration(ctx context.Context, params *protocol.DeclarationParams) ([]protocol.Location, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) DidChangeConfiguration(ctx context.Context, params *protocol.DidChangeConfigurationParams) error {
	return errors.New("unimplemented")
}

func (h Handler) DidChangeWatchedFiles(ctx context.Context, params *protocol.DidChangeWatchedFilesParams) error {
	return errors.New("unimplemented")
}

func (h Handler) DidChangeWorkspaceFolders(ctx context.Context, params *protocol.DidChangeWorkspaceFoldersParams) error {
	return errors.New("unimplemented")
}

func (h Handler) DocumentHighlight(ctx context.Context, params *protocol.DocumentHighlightParams) ([]protocol.DocumentHighlight, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) DocumentLink(ctx context.Context, params *protocol.DocumentLinkParams) ([]protocol.DocumentLink, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) DocumentLinkResolve(ctx context.Context, params *protocol.DocumentLink) (*protocol.DocumentLink, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) ExecuteCommand(ctx context.Context, params *protocol.ExecuteCommandParams) (interface{}, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) FoldingRanges(ctx context.Context, params *protocol.FoldingRangeParams) ([]protocol.FoldingRange, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) Formatting(ctx context.Context, params *protocol.DocumentFormattingParams) ([]protocol.TextEdit, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) Implementation(ctx context.Context, params *protocol.ImplementationParams) ([]protocol.Location, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) OnTypeFormatting(ctx context.Context, params *protocol.DocumentOnTypeFormattingParams) ([]protocol.TextEdit, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) PrepareRename(ctx context.Context, params *protocol.PrepareRenameParams) (*protocol.Range, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) RangeFormatting(ctx context.Context, params *protocol.DocumentRangeFormattingParams) ([]protocol.TextEdit, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) References(ctx context.Context, params *protocol.ReferenceParams) ([]protocol.Location, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) Rename(ctx context.Context, params *protocol.RenameParams) (*protocol.WorkspaceEdit, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) SignatureHelp(ctx context.Context, params *protocol.SignatureHelpParams) (*protocol.SignatureHelp, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) Symbols(ctx context.Context, params *protocol.WorkspaceSymbolParams) ([]protocol.SymbolInformation, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) TypeDefinition(ctx context.Context, params *protocol.TypeDefinitionParams) ([]protocol.Location, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) WillSave(ctx context.Context, params *protocol.WillSaveTextDocumentParams) error {
	return errors.New("unimplemented")
}

func (h Handler) WillSaveWaitUntil(ctx context.Context, params *protocol.WillSaveTextDocumentParams) ([]protocol.TextEdit, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) ShowDocument(ctx context.Context, params *protocol.ShowDocumentParams) (*protocol.ShowDocumentResult, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) WillCreateFiles(ctx context.Context, params *protocol.CreateFilesParams) (*protocol.WorkspaceEdit, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) DidCreateFiles(ctx context.Context, params *protocol.CreateFilesParams) error {
	return errors.New("unimplemented")
}

func (h Handler) WillRenameFiles(ctx context.Context, params *protocol.RenameFilesParams) (*protocol.WorkspaceEdit, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) DidRenameFiles(ctx context.Context, params *protocol.RenameFilesParams) error {
	return errors.New("unimplemented")
}

func (h Handler) WillDeleteFiles(ctx context.Context, params *protocol.DeleteFilesParams) (*protocol.WorkspaceEdit, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) DidDeleteFiles(ctx context.Context, params *protocol.DeleteFilesParams) error {
	return errors.New("unimplemented")
}

func (h Handler) CodeLensRefresh(ctx context.Context) error {
	return errors.New("unimplemented")
}

func (h Handler) PrepareCallHierarchy(ctx context.Context, params *protocol.CallHierarchyPrepareParams) ([]protocol.CallHierarchyItem, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) IncomingCalls(ctx context.Context, params *protocol.CallHierarchyIncomingCallsParams) ([]protocol.CallHierarchyIncomingCall, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) OutgoingCalls(ctx context.Context, params *protocol.CallHierarchyOutgoingCallsParams) ([]protocol.CallHierarchyOutgoingCall, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) SemanticTokensFull(ctx context.Context, params *protocol.SemanticTokensParams) (*protocol.SemanticTokens, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) SemanticTokensFullDelta(ctx context.Context, params *protocol.SemanticTokensDeltaParams) (interface{}, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) SemanticTokensRange(ctx context.Context, params *protocol.SemanticTokensRangeParams) (*protocol.SemanticTokens, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) SemanticTokensRefresh(ctx context.Context) error {
	return errors.New("unimplemented")
}

func (h Handler) LinkedEditingRange(ctx context.Context, params *protocol.LinkedEditingRangeParams) (*protocol.LinkedEditingRanges, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) Moniker(ctx context.Context, params *protocol.MonikerParams) ([]protocol.Moniker, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) Request(ctx context.Context, method string, params interface{}) (interface{}, error) {
	return nil, errors.New("unimplemented")
}
