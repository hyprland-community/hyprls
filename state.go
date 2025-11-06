package hyprls

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/hyprland-community/hyprls/parser"
	parser_data "github.com/hyprland-community/hyprls/parser/data"
	"go.lsp.dev/protocol"
	"go.uber.org/zap"
)

var logger *zap.Logger

var openedFiles = make(map[protocol.URI]string)
var ignores = []string{"hyprlock.conf", "hypridle.conf"} // should

type state struct {
}

func (h Handler) state(ctx context.Context) state {
	return ctx.Value("state").(state)
}

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

// TODO: possible to optimize
func currentLine(uri protocol.URI, position protocol.Position) (string, error) {
	contents, err := file(uri)
	if err != nil {
		return "", err
	}

	lines := strings.Split(contents, "\n")
	return lines[position.Line], nil
}

func isExclude(uri protocol.URI) bool {
	n := filepath.Base(uri.Filename())
	return slices.Contains(ignores, n)
}
