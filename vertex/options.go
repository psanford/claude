package vertex

import (
	"log/slog"
	"net/http"

	"golang.org/x/oauth2/google"
)

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

type regionOption struct {
	region string
}

func (o *regionOption) set(c *Client) {
	c.region = o.region
}

// Sets the region to connect to
func WithRegion(region string) Option {
	return &regionOption{
		region: region,
	}
}

type projectIDOption struct {
	projectID string
}

func (o *projectIDOption) set(c *Client) {
	c.projectID = o.projectID
}

// Sets the projectID to connect to
func WithProjectID(projectID string) Option {
	return &projectIDOption{
		projectID: projectID,
	}
}

type credOption struct {
	credentials *google.Credentials
}

func (o *credOption) set(c *Client) {
	c.credentials = o.credentials
}

// WithCredentials allows the user to pass in their own credentials
func WithCredentials(creds *google.Credentials) Option {
	return &credOption{
		credentials: creds,
	}
}
