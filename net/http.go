package net

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// HTTPClient describe the http client actions.
// In order to user the following struct properly should
// initialize the fields BaseURL and Http or to use the
// constructor NewHTTPClient.
type HTTPClient struct {
	// BaseURL defines the base url for all http requests.
	BaseURL string
	// Http is the Go HTTP client, abstracted for unit tests.
	Http Doer
}

// Doer abstraction of Go http client in order to be
// able to mocked for the unit tests.
type Doer interface {
	// Do default Go HTTP method for HTTP requests.
	Do(req *http.Request) (*http.Response, error)
}

// NewHTTPClient is responsible to return a HTTPClient object
// with all necessary parameter to use.
func NewHTTPClient(baseURL string, http Doer) HTTPClient {
	return HTTPClient{
		BaseURL: baseURL,
		Http:    http,
	}
}

// TCPAddressResponse describe the response of the server.
type TCPAddressResponse struct {
	Address string `json:"address"`
}

// GetTCPAddress is responsible to send an HTTP request
// in order to get the TCP address connection.
func (h HTTPClient) GetTCPAddress(token string) (string, error) {
	req, err := http.NewRequest("GET", h.BaseURL+"/tcp/connection", nil)
	if err != nil {
		return "", fmt.Errorf("failed to create http request error: %v", err)
	}

	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := h.Http.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed on request with error: %v", err)
	}
	defer resp.Body.Close()

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return "", fmt.Errorf("faield to read the body with error: %v", err)
	}

	var tcpAddress TCPAddressResponse
	err = json.Unmarshal(body, &tcpAddress)
	if err != nil {
		return "", fmt.Errorf("faield to unmarshal the body with error: %v", err)
	}

	return tcpAddress.Address, nil
}
