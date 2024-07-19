package claude

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestUnmarshalMessageResponse(t *testing.T) {
	respJSON := `{
  "content": [
    {
      "text": "Hi! My name is Claude.",
      "type": "text"
    }
  ],
  "id": "msg_013Zva2CMHLNnXjNJJKqJ2EF",
  "model": "claude-3-opus-20240229",
  "role": "assistant",
  "stop_reason": "end_turn",
  "stop_sequence": null,
  "type": "message",
  "usage": {
    "input_tokens": 10,
    "output_tokens": 25
  }
}`

	var mr MessageStart

	err := json.Unmarshal([]byte(respJSON), &mr)
	if err != nil {
		t.Fatal(err)
	}

	expect := MessageStart{
		ID:           "msg_013Zva2CMHLNnXjNJJKqJ2EF",
		Model:        "claude-3-opus-20240229",
		Role:         "assistant",
		StopReason:   "end_turn",
		StopSequence: nil,
		Type:         "message",
		Content: []TurnContent{
			TextContent("Hi! My name is Claude."),
		},
	}
	expect.Usage.InputTokens = 10
	expect.Usage.OutputTokens = 25

	if !cmp.Equal(mr, expect) {
		t.Fatalf(cmp.Diff(mr, expect))
	}

	respJSON = `{
  "type": "message_start",
  "message": {
	  "content": [
	    {
	      "text": "Hi! My name is Claude.",
	      "type": "text"
	    }
	  ],
	  "id": "msg_013Zva2CMHLNnXjNJJKqJ2EF",
	  "model": "claude-3-opus-20240229",
	  "role": "assistant",
	  "stop_reason": "end_turn",
	  "stop_sequence": null,
	  "type": "message",
	  "usage": {
	    "input_tokens": 10,
	    "output_tokens": 25
	  }
	}
}`

	mr = MessageStart{}
	err = json.Unmarshal([]byte(respJSON), &mr)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(mr, expect) {
		t.Fatalf(cmp.Diff(mr, expect))
	}
}

func TestMessageTurnUnmarshalJSON(t *testing.T) {
	mustDecodeB64 := func(s string) []byte {
		b, err := base64.StdEncoding.DecodeString(s)
		if err != nil {
			panic(err)
		}
		return b
	}

	tests := []struct {
		name     string
		input    string
		expected MessageTurn
		wantErr  bool
	}{
		{
			name: "Text content only",
			input: `{
				"role": "user",
				"content": [
					{"type": "text", "text": "Hello, world!"}
				]
			}`,
			expected: MessageTurn{
				Role: "user",
				Content: []TurnContent{
					&turnContentText{Typ: "text", Text: "Hello, world!"},
				},
			},
			wantErr: false,
		},
		{
			name: "Image content only",
			input: `{
				"role": "assistant",
				"content": [
					{"type": "image", "source": {"type": "base64", "media_type": "image/png", "data": "iVBORw0KGgo="}}
				]
			}`,
			expected: MessageTurn{
				Role: "assistant",
				Content: []TurnContent{
					&turnContentImage{
						Typ: "image",
						Source: struct {
							Type      string `json:"type"`
							MediaType string `json:"media_type"`
							Data      []byte `json:"data"`
						}{
							Type:      "base64",
							MediaType: "image/png",
							Data:      mustDecodeB64("iVBORw0KGgo="),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Mixed content",
			input: `{
				"role": "user",
				"content": [
					{"type": "text", "text": "Here's an image:"},
					{"type": "image", "source": {"type": "base64", "media_type": "image/jpeg", "data": "/9j/4AAQSkZJRg=="}},
					{"type": "text", "text": "What do you think?"}
				]
			}`,
			expected: MessageTurn{
				Role: "user",
				Content: []TurnContent{
					&turnContentText{Typ: "text", Text: "Here's an image:"},
					&turnContentImage{
						Typ: "image",
						Source: struct {
							Type      string `json:"type"`
							MediaType string `json:"media_type"`
							Data      []byte `json:"data"`
						}{
							Type:      "base64",
							MediaType: "image/jpeg",
							Data:      mustDecodeB64("/9j/4AAQSkZJRg=="),
						},
					},
					&turnContentText{Typ: "text", Text: "What do you think?"},
				},
			},
			wantErr: false,
		},
		{
			name: "Tool Use content",
			input: `{
				"role": "assistant",
				"content": [
					{"type": "tool_use", "id": "tool1", "name": "calculator", "input": {"operation": "add", "numbers": [1, 2, 3]}}
				]
			}`,
			expected: MessageTurn{
				Role: "assistant",
				Content: []TurnContent{
					&TurnContentToolUse{
						Typ:   "tool_use",
						ID:    "tool1",
						Name:  "calculator",
						Input: map[string]interface{}{"operation": "add", "numbers": []interface{}{float64(1), float64(2), float64(3)}},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Tool Result content",
			input: `{
				"role": "assistant",
				"content": [
					{"type": "tool_result", "tool_use_id": "tool1", "content": "The result is 6"}
				]
			}`,
			expected: MessageTurn{
				Role: "assistant",
				Content: []TurnContent{
					&turnContentToolResult{
						Typ:         "tool_result",
						ToolUseID:   "tool1",
						ToolContent: "The result is 6",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Unknown content type",
			input: `{
				"role": "user",
				"content": [
					{"type": "unknown", "data": "This should fail"}
				]
			}`,
			expected: MessageTurn{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got MessageTurn
			err := json.Unmarshal([]byte(tt.input), &got)

			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if diff := cmp.Diff(tt.expected, got); diff != "" {
					t.Errorf("UnmarshalJSON() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestUmarshalMessageStart(t *testing.T) {

}
