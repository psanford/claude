package vertex

import (
	"fmt"

	"github.com/psanford/claude"
)

type VertexModel string

const (
	Claude3Dot5SonnetV2 VertexModel = "claude-3-5-sonnet-v2@20241022"
	Claude3Dot5Haiku    VertexModel = "claude-3-5-haiku@20241022"
	Claude3Opus         VertexModel = "claude-3-opus@20240229"
	Claude3Dot5Sonnet   VertexModel = "claude-3-5-sonnet@20240620"
	Claude3Sonnet       VertexModel = "claude-3-sonnet@20240229"
	Claude3Haiku        VertexModel = "claude-3-haiku@20240307"
)

var Models = []VertexModel{
	Claude3Dot5SonnetV2,
	Claude3Dot5Haiku,
	Claude3Dot5Sonnet,
	Claude3Opus,
	Claude3Sonnet,
	Claude3Haiku,
}

func (m VertexModel) PrettyName() string {
	switch m {
	case Claude3Dot5SonnetV2:
		return "Claude 3.5 Sonnet V2"
	case Claude3Dot5Sonnet:
		return "Claude 3.5 Sonnet"
	case Claude3Dot5Haiku:
		return "Claude 3.5 Haiku"
	case Claude3Opus:
		return "Claude 3 Opus"
	case Claude3Sonnet:
		return "Claude 3 Sonnet"
	case Claude3Haiku:
		return "Claude 3 Haiku"
	default:
		return fmt.Sprintf("Unknown VertexModel<%s>", m)
	}
}

func ModelToVertexModel(m string) (VertexModel, error) {
	switch m {
	case claude.Claude3Dot5Sonnet2410, string(Claude3Dot5SonnetV2):
		return Claude3Dot5SonnetV2, nil
	case claude.Claude3Dot5Haiku, string(Claude3Dot5Haiku):
		return Claude3Dot5Haiku, nil
	case claude.Claude3Dot5Sonnet, string(Claude3Dot5Sonnet):
		return Claude3Dot5Sonnet, nil
	case claude.Claude3Opus, string(Claude3Opus):
		return Claude3Opus, nil
	case claude.Claude3Sonnet, string(Claude3Sonnet):
		return Claude3Sonnet, nil
	case claude.Claude3Haiku, string(Claude3Haiku):
		return Claude3Haiku, nil
	}

	return VertexModel(m), fmt.Errorf("Unknown model: %s", m)
}
