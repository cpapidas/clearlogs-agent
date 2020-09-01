package mock

import "net/http"

// HttpClient mock implementation
type HttpClient struct {
	SendLogsFunc func(token, log string) error
}

// SendLogs mock implementation
func (c *HttpClient) SendLogs(token, log string) error {
	if c.SendLogsFunc != nil {
		return c.SendLogsFunc(token, log)
	}
	return nil
}

// DoerMock mock implementation
type DoerMock struct {
	DoFunc func (req *http.Request) (*http.Response, error)
}

// Do mock implementation
func (c *DoerMock) Do(req *http.Request) (*http.Response, error) {
	if c.DoFunc != nil {
		return c.DoFunc(req)
	}
	return &http.Response{}, nil
}