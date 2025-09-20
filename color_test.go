package hyprls

import (
	"math/rand"
	"os"
	"testing"

	"go.lsp.dev/protocol"
	"go.uber.org/zap"
)

func init() {
	if logger == nil {
		if os.Getenv("HYPRLS_DEBUG") != "" {
			logger, _ = zap.NewDevelopmentConfig().Build()
		} else {
			logger, _ = zap.NewProductionConfig().Build()
		}
	}
}

func TestColorEncoding(t *testing.T) {
	for i := 0; i < 20; i++ {
		color := randomColor()
		colorstring := encodeColorLiteral(color)
		decoded := decodeColorLiteral(colorstring)
		encoded := encodeColorLiteral(decoded)
		if colorstring != encoded {
			t.Errorf("Color %d: %s != %s (was decoded to %v)", i, colorstring, encoded, decoded)
		}
	}
}

func TestColorDecoding(t *testing.T) {
	// red
	if !compareColorStructs(protocol.Color{Red: 1, Green: 0, Blue: 0, Alpha: 1}, decodeColorLiteral("rgb(ff0000)")) {
		t.Errorf("Color 'red' != {1, 0, 0, 1} (was decoded to %v)", decodeColorLiteral("red"))
	}

	// blue
	if !compareColorStructs(protocol.Color{Red: 0, Green: 0, Blue: 1, Alpha: 1}, decodeColorLiteral("rgba(0000ffff)")) {
		t.Errorf("Color 'blue' != {0, 0, 1, 1} (was decoded to %v)", decodeColorLiteral("blue"))
	}

	// green
	if !compareColorStructs(protocol.Color{Red: 0, Green: 1, Blue: 0, Alpha: 1}, decodeColorLiteral("0xff00ff00")) {
		t.Errorf("Color '0xff00ff00' != {0, 1, 0, 1} (was decoded to %v)", decodeColorLiteral("0xff00ff00"))
	}

	for i := 0; i < 20; i++ {
		color := randomColor()
		originalLogger := *logger
		logger = logger.WithOptions(zap.Fields(zap.Int("test color number", i)))
		encoded := encodeColorLiteral(color)
		decoded := decodeColorLiteral(encoded)
		if !compareColorStructs(color, decoded) {
			t.Errorf("Color %d: %#v != %#v (was encoded to %s)", i, color, decoded, encoded)
		}
		logger = &originalLogger
	}
}

func compareColorStructs(a, b protocol.Color) bool {
	delta := 0
	delta += int(a.Red*255 - b.Red*255)
	delta += int(a.Green*255 - b.Green*255)
	delta += int(a.Blue*255 - b.Blue*255)
	delta += int(a.Alpha*255 - b.Alpha*255)
	return delta == 0
}

func randomColor() protocol.Color {
	return protocol.Color{
		Red:   roundToThree(rand.Float64()),
		Green: roundToThree(rand.Float64()),
		Blue:  roundToThree(rand.Float64()),
		Alpha: roundToThree(rand.Float64()),
	}
}
