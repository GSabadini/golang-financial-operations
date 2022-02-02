package infrastructure

import (
	"io"
	"net/http"
)

// Client is the http wrapper for the application.
type Client struct {
	req *Request
}

// NewClient returns a configured Client.
func NewClient(r *Request) *Client {
	return &Client{r}
}

// Get executes a GET http request.
func (c *Client) Get(url string) (*http.Response, error) {
	return c.req.Do(http.MethodGet, url, nil)
}

// Post executes a POST http request.
func (c *Client) Post(url string, body io.Reader) (*http.Response, error) {
	return c.req.Do(http.MethodPost, url, body)
}
