package anthropic

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/psanford/claude"
	"github.com/psanford/claude/clientiface"
	"github.com/psanford/claude/internal/request"
	"github.com/psanford/claude/internal/responseparser"
)

var MessagesURL = "https://api.anthropic.com/v1/messages"

type Client struct {
	apiKey       string
	roundTripper http.RoundTripper
	debugLogger  *slog.Logger
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

func (c *Client) Message(ctx context.Context, req *claude.MessageRequest, options ...clientiface.Option) (claude.MessageResponse, error) {
	request.SetDefaults(req)
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
	headers.Add("anthropic-beta", "max-tokens-3-5-sonnet-2024-07-15")
	headers.Add("anthropic-beta", "output-128k-2025-02-19")
	headers.Add("x-api-key", c.apiKey)
	headers.Add("content-type", "application/json")
	httpReq.Header = headers

	client := c.httpClient()

	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}

	return responseparser.HandleResponse(ctx, resp, c.debugLogger)
}

func (c *Client) httpClient() *http.Client {
	if c.roundTripper == nil {
		return http.DefaultClient
	}
	return &http.Client{
		Transport: c.roundTripper,
	}
}
