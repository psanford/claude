package responseparser

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"net/http"

	"github.com/psanford/claude"
)

func HandleResponse(ctx context.Context, resp *http.Response, debugLogger *slog.Logger) (claude.MessageResponse, error) {
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
		return handleSSE(ctx, resp, debugLogger)
	} else if mediatype == "application/json" {
		return handleNonStreamingResponse(ctx, resp, debugLogger)
	} else {
		return nil, fmt.Errorf("unexpected response content-type: %s", contentType)
	}
}

func handleSSE(ctx context.Context, resp *http.Response, debugLogger *slog.Logger) (claude.MessageResponse, error) {
	eventsCh := decodeSSE(ctx, resp.Body)

	ch := make(chan claude.MessageEvent)
	meta := messageResponse{
		responses:    ch,
		httpResponse: resp,
	}

	go func() {
		defer close(ch)

		for evt := range eventsCh {
			if debugLogger != nil && debugLogger.Enabled(ctx, slog.LevelDebug) {
				debugLogger.Debug("sse event", "event", evt)
			}

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
			case "error":
				innerMsg = &claude.ClaudeError{}
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

func handleNonStreamingResponse(ctx context.Context, resp *http.Response, debugLogger *slog.Logger) (claude.MessageResponse, error) {
	d := json.NewDecoder(resp.Body)
	var msg claude.MessageStart
	err := d.Decode(&msg)
	if err != nil {
		return nil, err
	}
	if debugLogger != nil && debugLogger.Enabled(ctx, slog.LevelDebug) {
		debugLogger.Debug("response", "message_start", msg)
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
