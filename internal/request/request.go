package request

import "github.com/psanford/claude"

func SetDefaults(req *claude.MessageRequest) {
	if req.MaxTokens < 1 {
		req.MaxTokens = 4096
	}
}
