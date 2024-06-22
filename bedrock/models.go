package bedrock

import (
	"fmt"

	"github.com/psanford/claude"
)

// BedrockModel is the field set in the bedrock request in InvokeModel
// This is different than the model set in MessageRequest or TextCompletion
type BedrockModel string

const (
	Claude3Dot5Sonnet  BedrockModel = "anthropic.claude-3-5-sonnet-20240620-v1:0"
	Claude3Opus        BedrockModel = "anthropic.claude-3-opus-20240229-v1:0"
	Claude3Sonnet      BedrockModel = "anthropic.claude-3-sonnet-20240229-v1:0"
	Claude3Haiku       BedrockModel = "anthropic.claude-3-haiku-20240307-v1:0"
	Claude2Dot1        BedrockModel = "anthropic.claude-v2:1"
	Clause2Dot0        BedrockModel = "anthropic.claude-v2"
	Claude1Dot2Instant BedrockModel = "anthropic.claude-instant-v1"
)

var Models = []BedrockModel{
	Claude3Dot5Sonnet,
	Claude1Dot2Instant,
	Clause2Dot0,
	Claude2Dot1,
	Claude3Haiku,
	Claude3Sonnet,
	Claude3Opus,
}

func (m BedrockModel) PrettyName() string {
	switch m {
	case Claude3Dot5Sonnet:
		return "Claude 3.5 Sonnet"
	case Claude3Opus:
		return "Claude 3 Opus"
	case Claude3Sonnet:
		return "Claude 3 Sonnet"
	case Claude3Haiku:
		return "Claude 3 Haiku"
	case Claude2Dot1:
		return "Claude 2.1"
	case Clause2Dot0:
		return "Claude 2.0"
	case Claude1Dot2Instant:
		return "Claude 1.2 Instant"
	default:
		return fmt.Sprintf("Unknown BedrockModel<%s>", m)
	}
}

func ModelToBedrockModel(m string) (BedrockModel, error) {
	switch m {
	case claude.Claude3Dot5Sonnet, string(Claude3Dot5Sonnet):
		return Claude3Dot5Sonnet, nil
	case claude.Claude3Opus, string(Claude3Opus):
		return Claude3Opus, nil
	case claude.Claude3Sonnet, string(Claude3Sonnet):
		return Claude3Sonnet, nil
	case claude.Claude3Haiku, string(Claude3Haiku):
		return Claude3Haiku, nil
	case claude.Claude2Dot1, string(Claude2Dot1):
		return Claude2Dot1, nil
	case claude.Clause2Dot0, string(Clause2Dot0):
		return Clause2Dot0, nil
	case claude.Claude1Dot2Instant, string(Claude1Dot2Instant):
		return Claude1Dot2Instant, nil
	}

	return BedrockModel(m), fmt.Errorf("Unknown model: %s", m)

}
