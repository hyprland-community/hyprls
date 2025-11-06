package hyprls

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/hyprland-community/hyprls/parser"
	parser_data "github.com/hyprland-community/hyprls/parser/data"
	"go.lsp.dev/protocol"
	"go.uber.org/zap"
)

func (h Handler) Completion(ctx context.Context, params *protocol.CompletionParams) (*protocol.CompletionList, error) {
	if isExclude(params.TextDocument.URI) {
		return nil, nil
	}
	line, err := currentLine(params.TextDocument.URI, params.Position)
	if err != nil {
		return nil, nil
	}

	file, err := parse(params.TextDocument.URI)
	if err != nil {
		return nil, nil
	}

	sec := currentSection(file, params.Position)
	if sec == nil {
		sec = &parser.Section{}
	}

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
		if params.Position.Character <= 0 {
			return nil, nil
		}
		if !characterBeforeCursorIsDollarSign && !unicode.IsSpace(rune(line[params.Position.Character-1])) {
			return nil, nil
		}

		// Prevent duplicate dollar signs upon completion accept
		var textEditRange protocol.Range
		if characterBeforeCursorIsDollarSign {
			textEditRange = protocol.Range{
				Start: protocol.Position{Line: params.Position.Line, Character: params.Position.Character - 1},
				End:   protocol.Position{Line: params.Position.Line, Character: params.Position.Character},
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
