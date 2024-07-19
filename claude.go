package claude

import (
	"encoding/json"
	"fmt"
	"strings"
)

// TextCompletion represents the request to the legacy text completions api.
// This is deprecated. You should use the messages API vis MessageRequest instead.
// See https://docs.anthropic.com/claude/reference/complete_post for details
type TextCompletion struct {
	// The model that will complete your prompt.
	// Required field except for AWS Bedrock where it must be empty.
	Model string `json:"model,omitempty"`
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

// MessageRequest is a request struct for the messages API.
// See https://docs.anthropic.com/claude/reference/messages_post for details
type MessageRequest struct {
	// The model that will complete your prompt.
	// Required field except for AWS Bedrock where it must be empty.
	Model string `json:"model,omitempty"`
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
	// Recommended for advanced use cases only. You usually only need to use temperature.
	TopK *int `json:"top_k,omitempty"`
	// AnthropicVersion is used for AWS Bedrock and GCP Vertex.
	// The client implementations in this library will set this for you so you can leave it blank.
	AnthropicVersion string `json:"anthropic_version,omitempty"`
	// How the model should use the provided tools.
	ToolChoice *ToolChoice `json:"tool_choice,omitempty"`
	// Definitions of tools that the model may use.
	Tools []Tool `json:"tools,omitempty"`
}

// Tool defines a tool that the model may use.
type Tool struct {
	// Name of the tool.
	Name string `json:"name"`
	// Optional description of the tool.
	Description string `json:"description,omitempty"`
	// JSON schema for the tool input shape that the model will produce in tool_use output content blocks.
	InputSchema any `json:"input_schema"`
}

// ToolChoice defines how the model should use the provided tools.
type ToolChoice struct {
	// Specifies that the model should use a specific tool.
	Tool string `json:"tool,omitempty"`
	// Specifies that the model should use any available tool.
	Any bool `json:"any,omitempty"`
	// Specifies that the model should decide which tool to use.
	Auto bool `json:"auto,omitempty"`
}

type MessageResponse interface {
	Responses() <-chan MessageEvent
}

type MessageStart struct {
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

func (c *MessageStart) Text() string {
	text := make([]string, len(c.Content))
	for i, content := range c.Content {
		text[i] = content.TextContent()
	}
	return strings.Join(text, "")
}

func (m *MessageStart) UnmarshalJSON(b []byte) error {
	type ConcreteResponse struct {
		ID           string             `json:"id"`
		Type         string             `json:"type"`
		Role         string             `json:"role"`
		Content      []*turnContentText `json:"content"`
		Model        string             `json:"model"`
		StopReason   string             `json:"stop_reason"`
		StopSequence *string            `json:"stop_sequence"`
		Usage        struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
		} `json:"usage"`
	}

	type hackyBimodalResponse struct {
		*ConcreteResponse
		InnerMessage *ConcreteResponse `json:"message"`
	}

	var hack hackyBimodalResponse
	err := json.Unmarshal(b, &hack)
	if err != nil {
		return err
	}

	var c ConcreteResponse
	if hack.InnerMessage != nil {
		c = *hack.InnerMessage
	} else {
		c = *hack.ConcreteResponse
	}

	m.ID = c.ID
	m.Type = c.Type
	m.Role = c.Role
	m.Content = make([]TurnContent, len(c.Content))
	for i, c := range c.Content {
		m.Content[i] = c
	}
	m.Model = c.Model
	m.StopReason = c.StopReason
	m.StopSequence = c.StopSequence
	m.Usage = c.Usage

	return nil
}

type MessageTurn struct {
	Role    string        `json:"role"`
	Content []TurnContent `json:"content"`
}

func (m *MessageTurn) UnmarshalJSON(data []byte) error {
	var raw struct {
		Role    string            `json:"role"`
		Content []json.RawMessage `json:"content"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	m.Role = raw.Role
	m.Content = make([]TurnContent, len(raw.Content))

	for i, rawContent := range raw.Content {
		var contentType struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(rawContent, &contentType); err != nil {
			return err
		}

		switch contentType.Type {
		case TurnText:
			var textContent turnContentText
			if err := json.Unmarshal(rawContent, &textContent); err != nil {
				return err
			}
			m.Content[i] = &textContent

		case TurnImage:
			var imageContent turnContentImage
			if err := json.Unmarshal(rawContent, &imageContent); err != nil {
				return err
			}
			m.Content[i] = &imageContent

		case TurnToolUse:
			var toolUse TurnContentToolUse
			if err := json.Unmarshal(rawContent, &toolUse); err != nil {
				return err
			}
			m.Content[i] = &toolUse

		case TurnToolResult:
			var toolResult turnContentToolResult
			if err := json.Unmarshal(rawContent, &toolResult); err != nil {
				return err
			}
			m.Content[i] = &toolResult

		default:
			return fmt.Errorf("unknown content type: %s", contentType.Type)
		}
	}

	return nil
}

type TurnContent interface {
	Type() string
	TextContent() string
}

const (
	TurnText       = "text"
	TurnImage      = "image"
	TurnToolUse    = "tool_use"
	TurnToolResult = "tool_result"
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
	return TurnImage
}

func (t *turnContentImage) TextContent() string {
	return ""
}

type TurnContentToolUse struct {
	Typ   string      `json:"type"`
	ID    string      `json:"id"`
	Name  string      `json:"name"`
	Input interface{} `json:"input"`
}

func (t *TurnContentToolUse) Type() string {
	return TurnToolUse
}

func (t *TurnContentToolUse) TextContent() string {
	return ""
}

func ToolResultContent(toolUseID, content string) TurnContent {
	return &turnContentToolResult{
		Typ:         TurnToolResult,
		ToolUseID:   toolUseID,
		ToolContent: content,
	}
}

type turnContentToolResult struct {
	Typ         string `json:"type"`
	ToolUseID   string `json:"tool_use_id"`
	ToolContent string `json:"content"`
}

func (t *turnContentToolResult) Type() string {
	return TurnToolResult
}

func (t *turnContentToolResult) TextContent() string {
	return t.ToolContent
}

type MessageEvent struct {
	Type string
	Data MessageContent
}

type MessageContent interface {
	Text() string
}

type ContentBlockStart struct {
	ContentBlock struct {
		Text string `json:"text"`
		Type string `json:"type"`
		Name string `json:"name"`
		ID   string `json:"id"`
	} `json:"content_block"`
	Index int `json:"index"`
}

func (c *ContentBlockStart) Text() string {
	return c.ContentBlock.Text
}

type MessagePing struct {
}

func (c *MessagePing) Text() string {
	return ""
}

type ContentBlockDelta struct {
	Delta struct {
		Text        string `json:"text"`
		PartialJson string `json:"partial_json"`
		Type        string `json:"type"`
	} `json:"delta"`
	Index int64 `json:"index"`
}

func (c *ContentBlockDelta) Text() string {
	return c.Delta.Text
}

type ContentBlockStop struct {
	Index int64 `json:"index"`
}

func (c *ContentBlockStop) Text() string {
	return ""
}

type MessageDelta struct {
	Delta struct {
		StopReason   string  `json:"stop_reason"`
		StopSequence *string `json:"stop_sequence"`
	} `json:"delta"`
	Usage struct {
		OutputTokens int64 `json:"output_tokens"`
	} `json:"usage"`
}

func (c *MessageDelta) Text() string {
	return ""
}

type MessageStop struct {
}

func (c *MessageStop) Text() string {
	return ""
}

type ClaudeError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func (c *ClaudeError) Text() string {
	return ""
}

func (c ClaudeError) Error() string {
	return fmt.Sprintf("%s: %s", c.Type, c.Message)
}

type ClientError struct {
	error error
}

func NewClientError(err error) *ClientError {
	return &ClientError{
		error: err,
	}
}

func (c *ClientError) Error() string {
	return c.error.Error()
}

func (c *ClientError) Text() string {
	return ""
}
