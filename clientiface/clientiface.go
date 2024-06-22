package clientiface

import (
	"context"

	"github.com/psanford/claude"
)

type Client interface {
	Message(ctx context.Context, req *claude.MessageRequest) (claude.MessageResponse, error)
}
