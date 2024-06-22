# Go client library for Anthropic's Claude

This is an unofficial client library for Anthropic's Claude.
This project is not affiliated with Anthropic PBC.

## Package Layout

The `github.com/psanford/claude` package contains the API request and response message definitions. These are shared across the different API providers (Anthropic, AWS/Bedrock, GCP/Vertex).

- `github.com/psanford/claude/anthropic` contains an API client for using Anthropic's API.
- `github.com/psanford/claude/bedrock` contains an API client for using Claude in AWS Bedrock.
- `github.com/psanford/claude/vertex` contains an API client for using Claude in GCP Vertex.


Examples:
- [examples/anthropic-claude-cli-demo](https://github.com/psanford/claude/blob/main/examples/anthropic-claude-cli-demo/anthropic_claude_cli.go) for using Anthropic's API
- [examples/bedrock-claude-cli-demo](https://github.com/psanford/claude/blob/main/examples/bedrock-claude-cli-demo/bedrock_claude_cli.go) for using Anthropic Models via AWS Bedrock (you must get granted access to the models before using this API)
- [examples/vertex-claude-cli-demo](https://github.com/psanford/claude/blob/main/examples/vertex-claude-cli-demo/vertex_claude_cli.go) for using Anthropic Models via Google GCP Vertex (you must get granted access to the models before using this API)

## Design

The goal of this package is to give a consistent client experience across the different model hosting providers and across streaming vs non-streaming responses.

The Bedrock API mostly takes requests in the same shape as Anthropic's first party API, but not fully. The model IDs are different and need to be passed to bedrock differently, for example.

This package converts requests in the shape of the Anthropic first party API to the correct form for Bedrock and Vertex.

Likewise, the streaming vs non-streaming APIs are similar but not exactly the same. This package unifies streaming and non-streaming into a single interface so you can have one code path that can handle either.


## Example

```
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/psanford/claude"
	"github.com/psanford/claude/anthropic"
	"github.com/psanford/claude/bedrock"
	"github.com/psanford/claude/clientiface"
)

var apiProvider = flag.String("api", "anthropic", "API provider (anthropic|bedrock")
var stream = flag.Bool("stream", true, "Stream results")

func main() {
	flag.Parse()

	var client clientiface.Client
	if *apiProvider == "anthropic" {
		client = newAnthropicClient()
	} else if *apiProvider == "bedrock" {
		client = newBedrockClient()
	} else {
		log.Fatalf("Invalid api provider. Valid options are anthropic or bedrock")
	}

	err := makeRequestAndHandleResponse(client)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func newAnthropicClient() clientiface.Client {
	apiKey := os.Getenv("CLAUDE_API_KEY")
	if apiKey == "" {
		log.Fatal("CLAUDE_API_KEY environment variable is not set")
	}
	return anthropic.NewClient(apiKey)
}

func newBedrockClient() clientiface.Client {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("Unable to load AWS SDK config: %v", err)
	}
	bedrockSDK := bedrockruntime.NewFromConfig(cfg)
	return bedrock.NewClient(bedrockSDK)
}

func makeRequestAndHandleResponse(client clientiface.Client) error {
	req := &claude.MessageRequest{
		Model:     claude.Claude3Haiku,
		MaxTokens: 1000,
		Stream:    *stream, // Toggle streaming. No other code changes required
		Messages: []claude.MessageTurn{
			{
				Role: "user",
				Content: []claude.TurnContent{
					claude.TextContent("What are three interesting facts about Go programming?"),
				},
			},
		},
	}

	resp, err := client.Message(context.Background(), req)
	if err != nil {
		return fmt.Errorf("error calling Claude API: %w", err)
	}

	for event := range resp.Responses() {
		if err, isErr := event.Data.(error); isErr {
			return err
		}

		fmt.Print(event.Data.Text())
	}

	return nil
}
```
