package bedrock

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
	"github.com/psanford/claude"
	"github.com/psanford/claude/clientiface"
)

type Client struct {
	br *bedrockruntime.Client
}

var clientIfaceAssert = clientiface.Client(&Client{})

func NewClient(bedrockClient *bedrockruntime.Client, opts ...Option) *Client {
	c := &Client{
		br: bedrockClient,
	}
	for _, opt := range opts {
		opt.set(c)
	}
	return c
}

func (c *Client) Message(ctx context.Context, req *claude.MessageRequest, options ...clientiface.Option) (claude.MessageResponse, error) {
	if req.AnthropicVersion == "" {
		req.AnthropicVersion = "bedrock-2023-05-31"
	}

	bedrockModel, err := ModelToBedrockModel(req.Model)
	if err != nil {
		return nil, err
	}

	req.Model = "" // bedrock doesn't support this field here

	streaming := req.Stream
	req.Stream = false // bedrock doesn't support this field here

	jsonReq, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	if streaming {
		output, err := c.br.InvokeModelWithResponseStream(context.Background(), &bedrockruntime.InvokeModelWithResponseStreamInput{
			Body:        jsonReq,
			ModelId:     aws.String(string(bedrockModel)),
			ContentType: aws.String("application/json"),
		})

		if err != nil {
			return nil, err
		}

		return handleStreaming(ctx, output)
	} else {
		out, err := c.br.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
			Body:        jsonReq,
			ModelId:     aws.String(string(bedrockModel)),
			ContentType: aws.String("application/json"),
			Accept:      aws.String("application/json"),
		})

		if err != nil {
			return nil, err
		}

		var resp claude.MessageStart
		err = json.Unmarshal(out.Body, &resp)
		if err != nil {
			return nil, err
		}

		ch := make(chan claude.MessageEvent)
		meta := messageResponse{
			responses: ch,
		}

		evt := claude.MessageEvent{
			Type: resp.Type,
			Data: &resp,
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
}

func handleStreaming(ctx context.Context, output *bedrockruntime.InvokeModelWithResponseStreamOutput) (claude.MessageResponse, error) {

	ch := make(chan claude.MessageEvent)
	meta := messageResponse{
		responses: ch,
	}

	go func() {
		defer close(ch)

		for event := range output.GetStream().Events() {
			switch v := event.(type) {
			case *types.ResponseStreamMemberChunk:

				var msg claude.MessageEvent
				err := json.NewDecoder(bytes.NewReader(v.Value.Bytes)).Decode(&msg)
				if err != nil {
					msg = claude.MessageEvent{
						Type: "_client_error",
						Data: claude.NewClientError(fmt.Errorf("decode event json error: %w", err)),
					}
					select {
					case ch <- msg:
					case <-ctx.Done():
					}

					return
				}

				switch msg.Type {
				case "message_start":
					msg.Data = &claude.MessageStart{}
				case "ping":
					msg.Data = &claude.MessagePing{}
				case "content_block_start":
					msg.Data = &claude.ContentBlockStart{}
				case "content_block_delta":
					msg.Data = &claude.ContentBlockDelta{}
				case "content_block_stop":
					msg.Data = &claude.ContentBlockStop{}
				case "message_delta":
					msg.Data = &claude.MessageDelta{}
				case "message_stop":
					msg.Data = &claude.MessageStop{}
				default:
					msg = claude.MessageEvent{
						Type: "_client_error",
						Data: claude.NewClientError(fmt.Errorf("unknown event type: %s", msg.Type)),
					}
					select {
					case ch <- msg:
					case <-ctx.Done():
					}

					return
				}

				err = json.NewDecoder(bytes.NewReader(v.Value.Bytes)).Decode(&msg.Data)
				if err != nil {
					msg = claude.MessageEvent{
						Type: "_client_error",
						Data: claude.NewClientError(fmt.Errorf("decode event json error: %w", err)),
					}
					select {
					case ch <- msg:
					case <-ctx.Done():
					}

					return
				}

				select {
				case ch <- msg:
				case <-ctx.Done():
					return
				}
			case *types.UnknownUnionMember:
				msg := claude.MessageEvent{
					Type: "_client_error",
					Data: claude.NewClientError(fmt.Errorf("unknown bedrock tag: %s", v.Tag)),
				}
				select {
				case ch <- msg:
				case <-ctx.Done():
				}

				return
			default:
				msg := claude.MessageEvent{
					Type: "_client_error",
					Data: claude.NewClientError(fmt.Errorf("unknown bedrock event type: %T %+v", v, v)),
				}
				select {
				case ch <- msg:
				case <-ctx.Done():
				}

				return
			}
		}
	}()

	return &meta, nil
}

type messageResponse struct {
	responses <-chan claude.MessageEvent
}

func (m *messageResponse) Responses() <-chan claude.MessageEvent {
	return m.responses
}
