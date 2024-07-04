package bedrock

import "log/slog"

type Option interface {
	set(*Client)
}

type debugLoggerOption struct {
	l *slog.Logger
}

func (o *debugLoggerOption) set(c *Client) {
	c.debugLogger = o.l
}

func WithDebugLogger(l *slog.Logger) Option {
	return &debugLoggerOption{
		l: l,
	}
}
