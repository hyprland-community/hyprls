package hyprls

import (
	"github.com/ewen-lbh/hyprls/parser"
	"go.lsp.dev/protocol"
)

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
