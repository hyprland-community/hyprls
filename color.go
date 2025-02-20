package hyprls

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/hyprland-community/hyprls/parser"
	"github.com/mazznoer/csscolorparser"
	"go.lsp.dev/protocol"
	"go.uber.org/zap"
)

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

func decodeColorLiteral(raw string) protocol.Color {
	logger.Debug("decodeColorLiteral", zap.String("raw", raw))
	color, err := parser.ParseColor(raw)
	if err != nil {
		return protocol.Color{}
	}

	logger.Debug("decodeColorLiteral", zap.Any("color", color))

	return protocol.Color{
		Red:   roundToThree(float64(color.R) / 255.0),
		Alpha: roundToThree(float64(color.A) / 255.0),
		Blue:  roundToThree(float64(color.B) / 255.0),
		Green: roundToThree(float64(color.G) / 255.0),
	}
}

func encodeColorLiteral(color protocol.Color) string {
	logger.Debug("encodeColorLiteral", zap.Any("color", color))
	out := strings.TrimPrefix(csscolorparser.Color{
		R: roundToThree(color.Red),
		G: roundToThree(color.Green),
		B: roundToThree(color.Blue),
		A: roundToThree(color.Alpha),
	}.HexString(), "#")

	if color.Alpha == 1.0 {
		out = fmt.Sprintf("rgb(%s)", out)
	} else {
		out = fmt.Sprintf("rgba(%s)", out)
	}

	logger.Debug("encodeColorLiteral", zap.String("out", out))
	return out
}

func roundToThree(f float64) float64 {
	return math.Round(f*1_00) / 1_00
}
