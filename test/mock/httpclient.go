package mock

import "net/http"

// HttpClientMock mock implementation
type HttpClientMock struct {
	DoFunc func (req *http.Request) (*http.Response, error)
}

// Do mock implementation
func (c *HttpClientMock) Do(req *http.Request) (*http.Response, error) {
	if c.DoFunc != nil {
		return c.DoFunc(req)
	}
	return &http.Response{}, nil
}