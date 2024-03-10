package claude

import (
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

	var mr MessageResponse

	err := json.Unmarshal([]byte(respJSON), &mr)
	if err != nil {
		t.Fatal(err)
	}

	expect := MessageResponse{
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
}
