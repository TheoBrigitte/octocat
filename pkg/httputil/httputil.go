package httputil

import (
	"net/http"
)

// Client is a superset of http.Client
type Client struct {
	*http.Client
}

// New return a new Client.
func New() *Client {
	return &Client{
		&http.Client{},
	}
}

// NewWithTransport return a new client with the given http.Transport.
func NewWithTransport(tr *http.Transport) *Client {
	return &Client{
		&http.Client{Transport: tr},
	}
}

// Get issues a GET to the specified URL.
func (c Client) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

// Do sends an HTTP request and returns an HTTP response.
// It enforce header `Accept-Language: en-GB`.
func (c Client) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Accept-Language", "en-GB")

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
