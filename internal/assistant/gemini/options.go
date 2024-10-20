package gemini

import "github.com/lemoony/snipkit/internal/utils/httputil"

// Option configures a Manager.
type Option interface {
	apply(client *Client)
}

// optionFunc wraps a func so that it satisfies the Option interface.
type optionFunc func(client *Client)

func (f optionFunc) apply(client *Client) {
	f(client)
}

func WithConfig(config Config) Option {
	return optionFunc(func(client *Client) {
		client.config = config
	})
}

func WithHTTPClient(httpClient httputil.HTTPClient) Option {
	return optionFunc(func(c *Client) {
		c.httpClient = httpClient
	})
}
