package hyprls

import (
	"fmt"
	"math"
	"strings"

	"github.com/ewen-lbh/hyprlang-lsp/parser"
	"github.com/mazznoer/csscolorparser"
	"go.lsp.dev/protocol"
	"go.uber.org/zap"
)

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
