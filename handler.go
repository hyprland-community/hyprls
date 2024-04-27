package hyprls

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/MakeNowJust/heredoc"
	"github.com/ewen-lbh/hyprls/parser"
	parser_data "github.com/ewen-lbh/hyprls/parser/data"
	"go.lsp.dev/protocol"
	"go.uber.org/zap"
)

var logger *zap.Logger

var openedFiles = make(map[protocol.URI]string)

func parse(uri protocol.URI) (parser.Section, error) {
	contents, err := file(uri)
	if err != nil {
		return parser.Section{}, err
	}

	return parser.Parse(contents)
}

func currentSection(root parser.Section, position protocol.Position) *parser.Section {
	if !within(root.LSPRange(), position) {
		return nil
	}

	for _, section := range root.Subsections {
		sec := currentSection(section, position)
		if sec != nil {
			return sec
		}
	}

	return &root
}

func currentAssignment(root parser.Section, position protocol.Position) *parser_data.VariableDefinition {
	if !within(root.LSPRange(), position) {
		return nil
	}

	for _, assignment := range root.Assignments {
		if assignment.Position.Line == int(position.Line) {
			return parser_data.FindVariableDefinitionInSection(root.Name, assignment.Key)
		}
	}

	return nil
}

func within(rang protocol.Range, position protocol.Position) bool {
	if position.Line < rang.Start.Line || position.Line > rang.End.Line {
		return false
	}

	if position.Line == rang.Start.Line && position.Character < rang.Start.Character {
		return false
	}

	if position.Line == rang.End.Line && position.Character > rang.End.Character {
		return false
	}

	return true
}

func file(uri protocol.URI) (string, error) {
	if contents, ok := openedFiles[uri]; ok {
		return contents, nil
	}

	contents, err := os.ReadFile(uri.Filename())
	if err != nil {
		return "", err
	}

	openedFiles[uri] = string(contents)
	return string(contents), nil
}

func currentLine(uri protocol.URI, position protocol.Position) (string, error) {
	contents, err := file(uri)
	if err != nil {
		return "", err
	}

	lines := strings.Split(contents, "\n")
	return lines[position.Line], nil
}

type state struct {
}

type Handler struct {
	protocol.Server
	Logger *zap.Logger
}

func (h Handler) state(ctx context.Context) state {
	return ctx.Value("state").(state)
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
		},
		ServerInfo: &protocol.ServerInfo{
			Name:    "hyprls",
			Version: Version,
		},
	}, nil
}

func (h Handler) Definition(ctx context.Context, params *protocol.DefinitionParams) ([]protocol.Location, error) {
	return []protocol.Location{}, errors.New("unimplemented")
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

func (h Handler) ColorPresentation(ctx context.Context, params *protocol.ColorPresentationParams) ([]protocol.ColorPresentation, error) {
	logger.Debug("LSP:ColorPresentation", zap.Any("color", params.Color), zap.Any("range", params.Range))
	return []protocol.ColorPresentation{
		{
			Label: encodeColorLiteral(params.Color),
			TextEdit: &protocol.TextEdit{
				Range:   params.Range,
				NewText: encodeColorLiteral(params.Color),
			},
		},
	}, nil
}

func (h Handler) Completion(ctx context.Context, params *protocol.CompletionParams) (*protocol.CompletionList, error) {
	line, err := currentLine(params.TextDocument.URI, params.Position)
	file, err := parse(params.TextDocument.URI)
	if err != nil {
		return nil, nil
	}

	sec := currentSection(file, params.Position)

	cursorIsAfterEquals := err == nil && strings.Contains(line, "=") && strings.Index(line, "=") < int(params.Position.Character)

	// we are after the equals sign, suggest custom properties only
	if cursorIsAfterEquals {
		items := make([]protocol.CompletionItem, 0)

		cursorOrLineEnd := min(int(params.Position.Character), len(line)-1)
		characterBeforeCursorIsDollarSign := line[cursorOrLineEnd] == '$'

		// Don't propose custom variables if in the middle of typing a word
		// Only propose if a dollar sign was typed or is just before the cursor
		// Or we are after whitespace
		// Or we are in the middle of a color completion (typed a r, and key is a color or gradient)
		if !characterBeforeCursorIsDollarSign && !unicode.IsSpace(rune(line[params.Position.Character-1])) {
			return nil, nil
		}

		// Prevent duplicate dollar signs upon completion accept
		var textEditRange protocol.Range
		if characterBeforeCursorIsDollarSign {
			textEditRange = protocol.Range{
				protocol.Position{params.Position.Line, params.Position.Character - 1},
				protocol.Position{params.Position.Line, params.Position.Character},
			}
		} else {
			textEditRange = collapsedRange(params.Position)
		}

		file.WalkCustomVariables(func(v *parser.CustomVariable) {
			items = append(items, protocol.CompletionItem{
				Label: "$" + v.Key,
				Kind:  protocol.CompletionItemKindVariable,
				Documentation: protocol.MarkupContent{
					Kind:  protocol.PlainText,
					Value: v.ValueRaw,
				},
				TextEdit: &protocol.TextEdit{
					Range:   textEditRange,
					NewText: "$" + v.Key,
				},
			})
		})

		textedit := func(t string) *protocol.TextEdit {
			return &protocol.TextEdit{
				Range:   collapsedRange(params.Position),
				NewText: t,
			}
		}

		valueKind := parser.String
		if sec != nil {
			assignment := currentAssignment(*sec, params.Position)
			if assignment != nil {
				valueKind, err = parser.ValueKindFromString(assignment.ParserTypeString())
				if err != nil {
					valueKind = parser.String
				}
			}
		}

		logger.Debug("LSP:Completion", zap.Any("valueKind", valueKind))

		switch valueKind {
		case parser.Color, parser.Gradient:

			items = append(items, protocol.CompletionItem{
				Label:            "rgba(⋯)",
				Kind:             protocol.CompletionItemKindColor,
				InsertTextFormat: protocol.InsertTextFormatSnippet,
				Documentation:    "Define a color with an alpha channel of the form rgba(RRGGBBAA) in hexadecimal notation.",
				TextEdit:         textedit("rgba(${0:ffffffff})"),
			})

			items = append(items, protocol.CompletionItem{
				Label:            "rgb(⋯)",
				Kind:             protocol.CompletionItemKindColor,
				InsertTextFormat: protocol.InsertTextFormatSnippet,
				TextEdit:         textedit("rgb(${0:ffffff})"),
				Documentation:    "Define a color of the form rgb(RRGGBB) in hexadecimal notation.",
			})

			items = append(items, protocol.CompletionItem{
				Label:            "0xAARRGGBB",
				Kind:             protocol.CompletionItemKindColor,
				InsertTextFormat: protocol.InsertTextFormatSnippet,
				TextEdit:         textedit("0x${1:ffffffff}"),
				Deprecated:       true,
				Documentation:    "Define a color of the form 0xAARRGGBB in hexadecimal notation.",
			})
		case parser.Bool:
			items = append(items, protocol.CompletionItem{
				Label: "true",
				Kind:  protocol.CompletionItemKindValue,
			})
			items = append(items, protocol.CompletionItem{
				Label: "false",
				Kind:  protocol.CompletionItemKindValue,
			})
		case parser.Modmask:
			for keystring := range parser.ModKeyNames {
				items = append(items, protocol.CompletionItem{
					Label: keystring,
					Kind:  protocol.CompletionItemKindEnumMember,
				})
			}

		}

		return &protocol.CompletionList{
			Items: items,
		}, nil
	}

	availableVariables := make([]parser_data.VariableDefinition, 0)
	if sec != nil {
		secDef := parser_data.FindSectionDefinitionByName(sec.Name)
		if secDef != nil {
			availableVariables = append(availableVariables, secDef.Variables...)
		}
	}

	items := make([]protocol.CompletionItem, 0)
vars:
	for _, vardef := range availableVariables {
		// Don't suggest variables that are already defined
		for _, definedvar := range sec.Assignments {
			if vardef.Name == definedvar.Key {
				continue vars
			}
		}

		items = append(items, protocol.CompletionItem{
			Label: vardef.Name,
			Kind:  protocol.CompletionItemKindField,
			Documentation: protocol.MarkupContent{
				Kind:  protocol.Markdown,
				Value: fmt.Sprintf("Type: %s\n\n%s", vardef.Type, vardef.Description),
			},
		})
	}

	for _, kw := range parser_data.Keywords {
		items = append(items, protocol.CompletionItem{
			Label: kw.Name,
			Kind:  protocol.CompletionItemKindKeyword,
			Documentation: protocol.MarkupContent{
				Kind:  protocol.Markdown,
				Value: kw.Description,
			},
		})
	}

subsections:
	for _, subsections := range sec.Subsections {
		// Don't suggest subsections that are already defined
		for _, definedSubsection := range sec.Subsections {
			if subsections.Name == definedSubsection.Name {
				continue subsections
			}
		}

		items = append(items, protocol.CompletionItem{
			Label: subsections.Name,
			Kind:  protocol.CompletionItemKindModule,
			Documentation: protocol.MarkupContent{
				Kind:  protocol.Markdown,
				Value: fmt.Sprintf("Subsection of %s", sec.Name),
			},
		})
	}

	return &protocol.CompletionList{
		Items: items,
	}, nil
}

func (h Handler) CompletionResolve(ctx context.Context, params *protocol.CompletionItem) (*protocol.CompletionItem, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) Declaration(ctx context.Context, params *protocol.DeclarationParams) ([]protocol.Location, error) {
	return nil, errors.New("unimplemented")
}

func (h Handler) DidChange(ctx context.Context, params *protocol.DidChangeTextDocumentParams) error {
	logger.Debug("LSP:DidChange", zap.Any("params", params))
	openedFiles[params.TextDocument.URI] = params.ContentChanges[len(params.ContentChanges)-1].Text
	return nil
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

func (h Handler) DocumentColor(ctx context.Context, params *protocol.DocumentColorParams) ([]protocol.ColorInformation, error) {
	document, err := parse(params.TextDocument.URI)
	if err != nil {
		return []protocol.ColorInformation{}, fmt.Errorf("while parsing: %w", err)
	}
	colors := make([]protocol.ColorInformation, 0)
	document.WalkValues(func(a *parser.Assignment, v *parser.Value) {
		if v.Kind == parser.Gradient {
			for _, stop := range v.Gradient.Stops {
				colors = append(colors, protocol.ColorInformation{
					Color: stop.LSPColor(),
					Range: stop.LSPRange(),
				})
			}
			return
		}

		if v.Kind != parser.Color {
			return
		}

		colors = append(colors, protocol.ColorInformation{
			Color: v.LSPColor(),
			Range: v.LSPRange(),
		})
	})
	return colors, nil
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

func (h Handler) DocumentSymbol(ctx context.Context, params *protocol.DocumentSymbolParams) ([]interface{}, error) {
	document, err := parse(params.TextDocument.URI)
	if err != nil {
		return nil, fmt.Errorf("while parsing: %w", err)
	}
	symbols := make([]interface{}, 0)
	for _, symb := range gatherAllSymbols(document) {
		symbols = append(symbols, &symb)
	}
	return symbols, nil
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

func (h Handler) Hover(ctx context.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	line, err := currentLine(params.TextDocument.URI, params.Position)
	if err != nil {
		return nil, fmt.Errorf("while getting current line of file: %w", err)
	}

	if !strings.Contains(line, "=") {
		return nil, nil
	}

	// key is word before the equal sign. [0] is safe since we checked for "=" above
	key := strings.TrimSpace(strings.Split(line, "=")[0])

	indexOfFirstNonWhitespace := strings.IndexFunc(line, func(r rune) bool {
		return r != ' ' && r != '\t'
	})
	indexOfLastNonWhitespace := strings.LastIndexFunc(line, func(r rune) bool {
		return r != ' ' && r != '\t'
	}) + 1

	for _, section := range parser_data.Sections {
		if def := section.VariableDefinition(key); def != nil {
			return &protocol.Hover{
				Contents: protocol.MarkupContent{
					Kind: protocol.Markdown,
					Value: heredoc.Docf(`### %s: %s (%s)
						%s
						
						- Defaults to: %s
					`, strings.Join(section.Path, ":"), def.Name, def.Type, def.Description, def.PrettyDefault()),
				},
				Range: &protocol.Range{
					Start: protocol.Position{
						Line:      params.Position.Line,
						Character: uint32(indexOfFirstNonWhitespace),
					},
					End: protocol.Position{
						Line:      params.Position.Line,
						Character: uint32(indexOfLastNonWhitespace),
					},
				},
			}, nil
		} else if kw, found := parser_data.FindKeyword(key); found {
			flagsLine := ""
			if len(kw.Flags) > 0 {
				flagsLine = fmt.Sprintf("\n- Accepts the following flags: %s\n", strings.Join(kw.Flags, ", "))
			}
			return &protocol.Hover{
				Contents: protocol.MarkupContent{
					Kind:  protocol.Markdown,
					Value: fmt.Sprintf("### %s [[docs]](%s)%s\n%s", kw.Name, kw.DocumentationLink(), flagsLine, kw.Description),
				},
				Range: &protocol.Range{
					Start: protocol.Position{
						Line:      params.Position.Line,
						Character: uint32(indexOfFirstNonWhitespace),
					},
					End: protocol.Position{
						Line:      params.Position.Line,
						Character: uint32(indexOfLastNonWhitespace),
					},
				},
			}, nil
		}
	}

	return nil, nil
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
