# Go client library for Anthropic's Claude

This is an unofficial client library for Anthropic's Claude.
This project is not affiliated with Anthropic PBC.

## Package Layout

The `github.com/psanford/claude` package contains the API request and response message definitions. These are shared across the different API providers (Anthropic, AWS/Bedrock, GCP/Vertex).

`github.com/psanford/claude/anthropic` contains an API client for using Anthropic's API.

See [examples/anthropic-claude-cli-demo](https://github.com/psanford/claude/blob/main/examples/anthropic-claude-cli-demo/anthropic_claude_cli.go) for an example client using this API.
