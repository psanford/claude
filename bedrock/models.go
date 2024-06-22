package bedrock

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
