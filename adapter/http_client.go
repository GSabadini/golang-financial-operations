package adapter

import (
	"io"
	"net/http"
)

type (
	// HTTPClient is the http wrapper for the application
	HTTPClient interface {
		HTTPGetter
		HTTPPoster
	}

	// HTTPGetter holds fields and dependencies for executing an http GET request
	HTTPGetter interface {
		Get(url string) (*http.Response, error)
	}

	// HTTPPoster holds fields and dependencies for executing an http POST request
	HTTPPoster interface {
		Post(url string, body io.Reader) (*http.Response, error)
	}
)

type (
	stubHTTPClient struct {
		res *http.Response
		err error
	}
)

func (h stubHTTPClient) Get(_ string) (*http.Response, error) {
	return h.res, h.err
}

func (h stubHTTPClient) Post(_ string, _ io.Reader) (*http.Response, error) {
	return h.res, h.err
}
