package anthropic

import "net/http"

type Option interface {
	set(*Client)
}

type roundTripperOption struct {
	r http.RoundTripper
}

func (o *roundTripperOption) set(c *Client) {
	c.roundTripper = o.r
}

// Set a custom RoundTripper. This is useful if you wish to customize
// the http request before sending it.
func WithRoundTripper(r http.RoundTripper) Option {
	return &roundTripperOption{
		r: r,
	}
}
