package vertex

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/psanford/claude"
	"github.com/psanford/claude/clientiface"
	"github.com/psanford/claude/internal/request"
	"github.com/psanford/claude/internal/responseparser"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Client struct {
	projectID    string
	region       string
	credentials  *google.Credentials
	roundTripper http.RoundTripper
	debugLogger  *slog.Logger
}

var clientIfaceAssert = clientiface.Client(&Client{})

func NewClient(opts ...Option) *Client {
	c := &Client{
		region: os.Getenv("CLOUD_ML_REGION"), // default to env variable if not set in an option.
	}
	for _, opt := range opts {
		opt.set(c)
	}
	return c
}

func (c *Client) Message(ctx context.Context, req *claude.MessageRequest, options ...clientiface.Option) (claude.MessageResponse, error) {
	request.SetDefaults(req)

	if req.AnthropicVersion == "" {
		req.AnthropicVersion = "vertex-2023-10-16"
	}

	if c.region == "" {
		return nil, fmt.Errorf("region not set or automatically detected")
	}

	if c.projectID == "" {
		return nil, fmt.Errorf("projectID not set or automatically detected")
	}

	vertexModel, err := ModelToVertexModel(req.Model)
	if err != nil {
		return nil, err
	}
	req.Model = ""

	jsonReq, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	apiMethod := "rawPredict"
	if req.Stream {
		apiMethod = "streamRawPredict"
	}

	messageURL := fmt.Sprintf("https://%s-aiplatform.googleapis.com/v1/projects/%s/locations/%s/publishers/anthropic/models/%s:%s", c.region, c.projectID, c.region, vertexModel, apiMethod)

	httpReq, err := http.NewRequestWithContext(ctx, "POST", messageURL, bytes.NewReader(jsonReq))
	if err != nil {
		return nil, err
	}

	headers := make(http.Header)
	headers.Add("content-type", "application/json")
	httpReq.Header = headers

	client, err := c.httpClient(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}

	return responseparser.HandleResponse(ctx, resp)
}

func (c *Client) httpClient(ctx context.Context) (*http.Client, error) {
	if c.roundTripper != nil {
		return &http.Client{
			Transport: c.roundTripper,
		}, nil
	}

	var creds *google.Credentials
	var err error

	if c.credentials != nil {
		creds = c.credentials
	} else {
		creds, err = google.FindDefaultCredentials(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to find default credentials: %v", err)
		}
	}

	if c.projectID == "" {
		c.projectID = creds.ProjectID
	}

	return oauth2.NewClient(ctx, creds.TokenSource), nil
}
