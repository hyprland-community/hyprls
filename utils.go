package hyprls

import (
	"go.lsp.dev/protocol"
)

func collapsedRange(position protocol.Position) protocol.Range {
	return protocol.Range{
		Start: position,
		End:   position,
	}
}
