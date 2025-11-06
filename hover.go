package hyprls

import (
	"context"
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	parser_data "github.com/hyprland-community/hyprls/parser/data"
	"go.lsp.dev/protocol"
)

func (h Handler) Hover(ctx context.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	if isExclude(params.TextDocument.URI) {
		return nil, nil
	}
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
