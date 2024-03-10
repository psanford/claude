package claude

// TextCompletion represents the request to the legacy text completions api.
// This is deprecated. You should use the messages API vis MessageRequest instead.
// See https://docs.anthropic.com/claude/reference/complete_post for details
type TextCompletion struct {
	// The model that will complete your prompt.
	Model string `json:"model"`
	// The prompt that you want Claude to complete.
	// For proper response generation you will need to format your prompt using alternating
	// \n\nHuman: and \n\nAssistant: conversational turns.
	Prompt string `json:"prompt"`
	// The maximum number of tokens to generate before stopping.
	// Note that models may stop before reaching this maximum.
	// This parameter only specifies the absolute maximum number of tokens to generate.
	MaxTokensToSample int `json:"max_tokens_to_sample"`
	// Sequences that will cause the model to stop generating.
	// Models stop on "\n\nHuman:", and may include additional built-in stop sequences in the future.
	// By providing the stop_sequences parameter, you may include additional strings that will cause the model to stop generating.
	StopSequences []string `json:"stop_sequences,omitempty"`
	// Amount of randomness injected into the response.
	// Defaults to 1.0. Ranges from 0.0 to 1.0. Use temperature closer to 0.0 for analytical / multiple choice,
	// and closer to 1.0 for creative and generative tasks.
	// Note that even with temperature of 0.0, the results will not be fully deterministic.
	Temperature *float64 `json:"temperature,omitempty"`
	// Use nucleus sampling.
	// In nucleus sampling, we compute the cumulative distribution over all the options for each subsequent
	// token in decreasing probability order and cut it off once it reaches a particular probability specified by top_p.
	// You should either alter temperature or top_p, but not both.
	// Recommended for advanced use cases only. You usually only need to use temperature.
	TopP *float64 `json:"top_p,omitempty"`
	// Only sample from the top K options for each subsequent token.
	// Used to remove "long tail" low probability responses.
	// Recommended for advanced use cases only. You usually only need to use temperature
	TopK *int `json:"top_k,omitempty"`
	// An object describing metadata about the request.
	Metadata *RequestMetadata `json:"metadata"`
	// Whether to incrementally stream the response using server-sent events.
	Stream bool `json:"stream,omitempty"`
}

type RequestMetadata struct {
	// An external identifier for the user who is associated with the request.
	// This should be a uuid, hash value, or other opaque identifier.
	// Anthropic may use this id to help detect abuse.
	// Do not include any identifying information such as name, email address, or phone number.
	UserID string `json:"user_id,omitempty"`
}

type TextCompletionResponse struct {
	Type       string `json:"type"`
	ID         string `json:"id"`
	Completion string `json:"completion"`
	StopReason string `json:"stop_reason"`
	Model      string `json:"model"`
}

type ErrorResponse struct {
	Type  string `json:"type"`
	Error struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error"`
}

// MessageRequest is a request struct for the messages API.
// See https://docs.anthropic.com/claude/reference/messages_post for details
type MessageRequest struct {
	// The model that will complete your prompt.
	Model string `json:"model"`
	// Input messages.
	// Models are trained to operate on alternating user and assistant conversational turns.
	// When creating a new Message, you specify the prior conversational turns with the messages parameter,
	// and the model then generates the next Message in the conversation.
	// Each input message must be an object with a role and content. You can specify a single user-role message,
	// or you can include multiple user and assistant messages. The first message must always use the user role.
	// If the final message uses the assistant role, the response content will continue immediately from the
	// content in that message. This can be used to constrain part of the model's response.
	Messages []MessageTurn `json:"messages"`
	// System prompt.
	// A system prompt is a way of providing context and instructions to Claude, such as specifying a particular goal or role.
	System string `json:"system,omitempty"`
	// The maximum number of tokens to generate before stopping.
	// Note that models may stop before reaching this maximum.
	// This parameter only specifies the absolute maximum number of tokens to generate.
	// Different models have different maximum values for this parameter.
	MaxTokens int              `json:"max_tokens"`
	Metadata  *RequestMetadata `json:"metadata,omitempty"`
	// Custom text sequences that will cause the model to stop generating.
	// Models will normally stop when they have naturally completed their turn,
	// which will result in a response stop_reason of "end_turn".
	// If you want the model to stop generating when it encounters custom strings of text,
	// you can use the stop_sequences parameter. If the model encounters one of the custom sequences,
	// the response stop_reason value will be "stop_sequence" and the response stop_sequence value will contain the matched stop sequence.
	StopSequences []string `json:"stop_sequences,omitempty"`
	// Whether to incrementally stream the response using server-sent events.
	Stream bool `json:"stream,omitempty"`
	// Amount of randomness injected into the response.
	// Defaults to 1.0. Ranges from 0.0 to 1.0. Use temperature closer to 0.0 for analytical / multiple choice, and closer to 1.0 for creative and generative tasks.
	// Note that even with temperature of 0.0, the results will not be fully deterministic.
	Temperature *float64 `json:"temperature,omitempty"`
	// Use nucleus sampling.
	// In nucleus sampling, we compute the cumulative distribution over all the options for each subsequent
	// token in decreasing probability order and cut it off once it reaches a particular probability specified by top_p.
	// You should either alter temperature or top_p, but not both.
	// Recommended for advanced use cases only. You usually only need to use temperature.
	TopP *float64 `json:"top_p,omitempty"`
	// Only sample from the top K options for each subsequent token.
	// Used to remove "long tail" low probability responses.
	// Recommended for advanced use cases only. You usually only need to use temperature
	TopK *int `json:"top_k,omitempty"`
	// AnthropicVersion is an AWS Bedrock specific field. For bedrock, it should be set to "bedrock-2023-05-31".
	// Otherwise it should be left blank.
	AnthropicVersion string `json:"anthropic_version,omitempty"`
}

type MessageResponse struct {
	ID           string        `json:"id"`
	Type         string        `json:"type"`
	Role         string        `json:"role"`
	Content      []TurnContent `json:"content"`
	Model        string        `json:"model"`
	StopReason   string        `json:"stop_reason"`
	StopSequence *string       `json:"stop_sequence"`
	Usage        struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}

type MessageTurn struct {
	Role    string        `json:"role"`
	Content []TurnContent `json:"content"`
}

type TurnContent interface {
	Type() string
	TextContent() string
}

const (
	TurnText  = "text"
	TurnImage = "image"
)

func TextContent(msg string) TurnContent {
	return &turnContentText{
		Typ:  TurnText,
		Text: msg,
	}
}

func (t *turnContentText) Type() string {
	return TurnText
}

func (t *turnContentText) TextContent() string {
	return t.Text
}

type turnContentText struct {
	Typ  string `json:"type"`
	Text string `json:"text"`
}

type turnContentImage struct {
	Typ    string `json:"type"`
	Source struct {
		Type      string `json:"type"`
		MediaType string `json:"media_type"`
		Data      []byte `json:"data"`
	} `json:"source"`
}

func ImageContent(mediaType string, image []byte) TurnContent {
	i := turnContentImage{
		Typ: TurnImage,
	}
	i.Source.Type = "base64"
	i.Source.MediaType = mediaType
	i.Source.Data = image
	return &i
}

func (t *turnContentImage) Type() string {
	return TurnText
}

func (t *turnContentImage) TextContent() string {
	return ""
}
