package hyprls

import (
	"context"
	"fmt"

	"github.com/hyprland-community/hyprls/parser"
	"go.lsp.dev/protocol"
)

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

func gatherAllSymbols(root parser.Section) []protocol.DocumentSymbol {
	symbols := make([]protocol.DocumentSymbol, 0)
	for _, variable := range root.Assignments {
		symbols = append(symbols, protocol.DocumentSymbol{
			Name:           variable.Key,
			Kind:           variable.Value.Kind.LSPSymbol(),
			Detail:         variable.ValueRaw,
			Range:          collapsedRange(variable.Position.LSP()),
			SelectionRange: collapsedRange(variable.Position.LSP()),
		})
	}
	for _, customVar := range root.Variables {
		symbols = append(symbols, protocol.DocumentSymbol{
			Name:           "$" + customVar.Key,
			Kind:           protocol.SymbolKindVariable,
			Detail:         customVar.ValueRaw,
			Range:          collapsedRange(customVar.Position.LSP()),
			SelectionRange: collapsedRange(customVar.Position.LSP()),
		})
	}
	for _, section := range root.Subsections {
		symbols = append(symbols, protocol.DocumentSymbol{
			Name:           section.Name,
			Kind:           protocol.SymbolKindNamespace,
			Range:          section.LSPRange(),
			SelectionRange: section.LSPRange(),
			Children:       gatherAllSymbols(section),
		})
	}
	return symbols
}
