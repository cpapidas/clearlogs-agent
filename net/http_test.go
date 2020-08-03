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
	want := "9.3.3.2:8876"
	mockHttp := &mock.HttpClientMock{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			resp := &http.Response{
				Body: ioutil.NopCloser(strings.NewReader(`{"address": "` + want + `"}`)),
			}
			return resp, nil
		},
	}
	got, err := NewHTTPClient("https://example.com", mockHttp).GetTCPAddress("token")
	if err != nil {
		t.Fatalf("expected error to be nil, but got: %v", err)
	}
	if want != got {
		t.Errorf("expected: %s, got: %s", want, got)
	}
}

func TestHTTPClientGetTCPAddressShouldReturnAnErrorForInvalidRequest(t *testing.T) {
	want := ""
	mockHttp := &mock.HttpClientMock{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("error")
		},
	}
	got, err := NewHTTPClient("https://example.com", mockHttp).GetTCPAddress("token")
	if err == nil {
		t.Fatal("expected error to not be nil")
	}
	if want != got {
		t.Errorf("expected: %s, got: %s", want, got)
	}
}
