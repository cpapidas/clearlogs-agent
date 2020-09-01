package net

import (
	"errors"
	"github.com/cpapidas/clagent/test/mock"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestHTTPClientGetTCPAddressShouldReturnTheAddressForAValidHTTPRequest(t *testing.T) {
	mockHttp := &mock.DoerMock{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			resp := &http.Response{
				Body: ioutil.NopCloser(strings.NewReader(`{"success": "ok"}`)),
			}
			return resp, nil
		},
	}
	err := NewHTTPClient("https://example.com", mockHttp).SendLogs("token", "payload")
	if err != nil {
		t.Fatalf("expected error to be nil, but got: %v", err)
	}
}

func TestHTTPClientGetTCPAddressShouldReturnAnErrorForInvalidRequest(t *testing.T) {
	mockHttp := &mock.DoerMock{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("error")
		},
	}
	err := NewHTTPClient("https://example.com", mockHttp).SendLogs("token", "payload")
	if err == nil {
		t.Fatal("expected error to not be nil")
	}
}
