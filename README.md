# Go client library for Anthropic's Claude

This is an unofficial client library for Anthropic's Claude.
This project is not affiliated with Anthropic PBC.

## Package Layout

The `github.com/psanford/claude` package contains the API request and response message definitions. These are shared across the different API providers (Anthropic, AWS/Bedrock, GCP/Vertex).

`github.com/psanford/claude/anthropic` contains an API client for using Anthropic's API.

Examples:
- [examples/anthropic-claude-cli-demo](https://github.com/psanford/claude/blob/main/examples/anthropic-claude-cli-demo/anthropic_claude_cli.go) for using Anthropic's API
- [examples/bedrock-claude-cli-demo](https://github.com/psanford/claude/blob/main/examples/bedrock-claude-cli-demo/bedrock_claude_cli.go) for using Anthropic Models via AWS Bedrock (you must get granted access to the models before using this API)

## Design

The goal of this package is to give a consistent client experience across the different model hosting providers and across streaming vs non-streaming responses.

The Bedrock API mostly takes requests in the same shape as Anthropic's first party API, but not fully. The model IDs are different and need to be passed to bedrock differently, for example.

This package converts requests in the shape of the Anthropic first party API to the correct form for Bedrock (and eventually Vertex).

Likewise, the streaming vs non-streaming APIs are similar but not exactly the same. This package unifies streaming and non-streaming into a single interface so you can have one code path that can handle either.
