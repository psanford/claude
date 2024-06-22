package anthropic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"

	"github.com/psanford/claude"
	"github.com/psanford/claude/clientiface"
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

	if resp.StatusCode != 200 {
		r := io.LimitReader(resp.Body, 1<<13)
		body, err := io.ReadAll(r)
		if err != nil {
			return nil, fmt.Errorf("%d error response: %s", resp.StatusCode, body)
		}
		var ew errWrapper
		err = json.Unmarshal(body, &ew)
		if err != nil {
			return nil, fmt.Errorf("%d error response: %s", resp.StatusCode, body)
		}
		return nil, ew.Error
	}

	contentType := resp.Header.Get("Content-Type")
	mediatype, _, _ := mime.ParseMediaType(contentType)
	if mediatype == "text/event-stream" {
		return c.handleSSE(ctx, resp)
	} else if mediatype == "application/json" {
		return c.handleNonStreamingResponse(ctx, resp)
	} else {
		return nil, fmt.Errorf("unexpected response content-type: %s", contentType)
	}
}

func (c *Client) handleSSE(ctx context.Context, resp *http.Response) (claude.MessageResponse, error) {
	eventsCh := decodeSSE(ctx, resp.Body)

	ch := make(chan claude.MessageEvent)
	meta := messageResponse{
		responses:    ch,
		httpResponse: resp,
	}

	go func() {
		defer close(ch)

		for evt := range eventsCh {
			var msg claude.MessageEvent
			if evt.Error != nil {
				msg = claude.MessageEvent{
					Type: "_client_error",
					Data: claude.NewClientError(evt.Error),
				}
				select {
				case ch <- msg:
				case <-ctx.Done():
				}

				return
			}

			msg.Type = evt.Name

			var innerMsg claude.MessageContent

			switch evt.Name {
			case "message_start":
				innerMsg = &claude.MessageStart{}
			case "ping":
				innerMsg = &claude.MessagePing{}
			case "content_block_start":
				innerMsg = &claude.ContentBlockStart{}
			case "content_block_delta":
				innerMsg = &claude.ContentBlockDelta{}
			case "content_block_stop":
				innerMsg = &claude.ContentBlockStop{}
			case "message_delta":
				innerMsg = &claude.MessageDelta{}
			case "message_stop":
				innerMsg = &claude.MessageStop{}
			default:
				msg = claude.MessageEvent{
					Type: "_client_error",
					Data: claude.NewClientError(fmt.Errorf("unknown event type: %s", evt.Name)),
				}
				select {
				case ch <- msg:
				case <-ctx.Done():
				}

				return
			}

			err := json.Unmarshal([]byte(evt.Data), innerMsg)
			if err != nil {
				msg = claude.MessageEvent{
					Type: "_client_error",
					Data: claude.NewClientError(fmt.Errorf("parse event err: %w", err)),
				}
				select {
				case ch <- msg:
				case <-ctx.Done():
				}

				return
			}

			msg.Data = innerMsg
			select {
			case ch <- msg:
			case <-ctx.Done():
				return
			}
		}
	}()

	return &meta, nil
}

func (c *Client) handleNonStreamingResponse(ctx context.Context, resp *http.Response) (claude.MessageResponse, error) {
	d := json.NewDecoder(resp.Body)
	var msg claude.MessageStart
	err := d.Decode(&msg)
	if err != nil {
		return nil, err
	}

	ch := make(chan claude.MessageEvent)
	meta := messageResponse{
		responses:    ch,
		httpResponse: resp,
	}

	evt := claude.MessageEvent{
		Type: msg.Type,
		Data: &msg,
	}
	go func() {
		select {
		case ch <- evt:
		case <-ctx.Done():
		}

		close(ch)
	}()

	return &meta, nil
}

type messageResponse struct {
	responses    <-chan claude.MessageEvent
	httpResponse *http.Response
}

func (m *messageResponse) Responses() <-chan claude.MessageEvent {
	return m.responses
}

func (m *messageResponse) HTTPResponse() *http.Response {
	return m.httpResponse
}

func (c *Client) httpClient() *http.Client {
	if c.roundTripper == nil {
		return http.DefaultClient
	}
	return &http.Client{
		Transport: c.roundTripper,
	}
}

type Error struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

type errWrapper struct {
	Type  string `json:"type"`
	Error Error  `json:"error"`
}
