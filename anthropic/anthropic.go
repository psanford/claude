package anthropic

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/psanford/claude"
	"github.com/psanford/claude/clientiface"
	"github.com/psanford/claude/internal/responseparser"
)

var MessagesURL = "https://api.anthropic.com/v1/messages"

type Client struct {
	apiKey       string
	roundTripper http.RoundTripper
}

var clientIfaceAssert = clientiface.Client(&Client{})

func NewClient(apiKey string, opts ...Option) *Client {
	c := &Client{
		apiKey: apiKey,
	}
	for _, opt := range opts {
		opt.set(c)
	}
	return c
}

func (c *Client) Message(ctx context.Context, req *claude.MessageRequest) (claude.MessageResponse, error) {
	jsonReq, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", MessagesURL, bytes.NewReader(jsonReq))
	if err != nil {
		return nil, err
	}

	headers := make(http.Header)
	headers.Add("anthropic-version", "2023-06-01")
	headers.Add("x-api-key", c.apiKey)
	headers.Add("content-type", "application/json")
	httpReq.Header = headers

	client := c.httpClient()

	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}

	return responseparser.HandleResponse(ctx, resp)
}

func (c *Client) httpClient() *http.Client {
	if c.roundTripper == nil {
		return http.DefaultClient
	}
	return &http.Client{
		Transport: c.roundTripper,
	}
}
