package helper

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPService defines an interface for making http requests
type HTTPService interface {
	Get(string) (*http.Response, error)
	GetWithHeaders(url string, params map[string]string, headers map[string]string) (*http.Response, error)
	PutWithHeaders(url string, contentType string, headers map[string]string, body io.Reader) (*http.Response, error)
	PatchWithHeaders(url string, contentType string, headers map[string]string, body io.Reader) (*http.Response, error)
	DeleteWithHeaders(url string, headers map[string]string) (*http.Response, error)
	Post(url string, contentType string, body io.Reader) (*http.Response, error)
	PostWithHeaders(url string, contentType string, headers map[string]string, body io.Reader) (*http.Response, error)
	PostEmptyWithHeaders(url string, contentType string, headers map[string]string) (*http.Response, error)
}

// HTTPServiceContainer returns an http service
type HTTPServiceContainer struct {
}

// NewHTTPService returns an http service
func NewHTTPService() HTTPService {
	service := &HTTPServiceContainer{}
	return service
}

// Get wraps http.Get to satisfy the HTTPService interface
func (hsc *HTTPServiceContainer) Get(url string) (*http.Response, error) {
	var netClient = &http.Client{
		Timeout: time.Second * 30,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()
	fmt.Println(req.URL.String())

	response, err := netClient.Do(req)
	if err != nil {
		return nil, err
	}

	return response, err
}

// GetWithHeaders wraps http.NewRequest with method GET and set request headers
func (hsc *HTTPServiceContainer) GetWithHeaders(url string, params, headers map[string]string) (*http.Response, error) {
	var netClient = &http.Client{
		Timeout: time.Second * 30,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	for key, value := range params {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()
	fmt.Println(req.URL.String())

	for headerKey, headerValue := range headers {
		req.Header.Add(headerKey, headerValue)
	}

	response, err := netClient.Do(req)
	if err != nil {
		return nil, err
	}

	return response, err
}

// PutWithHeaders wraps http.NewRequest with method PUT and set request headers
func (hsc *HTTPServiceContainer) PutWithHeaders(url string, contentType string, headers map[string]string, body io.Reader) (*http.Response, error) {
	var netClient = &http.Client{
		Timeout: time.Second * 30,
	}
	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)
	for headerKey, headerValue := range headers {
		req.Header.Add(headerKey, headerValue)
	}

	response, err := netClient.Do(req)
	return response, err
}

// PatchWithHeaders wraps http.NewRequest with method PATCH and set request headers
func (hsc *HTTPServiceContainer) PatchWithHeaders(url string, contentType string, headers map[string]string, body io.Reader) (*http.Response, error) {
	var netClient = &http.Client{
		Timeout: time.Second * 30,
	}
	req, err := http.NewRequest("PATCH", url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)
	for headerKey, headerValue := range headers {
		req.Header.Add(headerKey, headerValue)
	}

	response, err := netClient.Do(req)
	return response, err
}

// DeleteWithHeaders wraps http.NewRequest with method DELETE
func (hsc *HTTPServiceContainer) DeleteWithHeaders(url string, headers map[string]string) (*http.Response, error) {
	var netClient = &http.Client{
		Timeout: time.Second * 30,
	}
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	for headerKey, headerValue := range headers {
		req.Header.Add(headerKey, headerValue)
	}

	response, err := netClient.Do(req)
	return response, err
}

// PostWithHeaders wraps http.NewRequest with method POST and set request headers
func (hsc *HTTPServiceContainer) PostWithHeaders(url string, contentType string, headers map[string]string, body io.Reader) (*http.Response, error) {
	var netClient = &http.Client{
		Timeout: time.Second * 30,
	}
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)
	for headerKey, headerValue := range headers {
		req.Header.Add(headerKey, headerValue)
	}

	response, err := netClient.Do(req)
	return response, err
}

// Post wraps http.Post to satisfy the HTTPService interface
func (hsc *HTTPServiceContainer) Post(url, contentType string, body io.Reader) (*http.Response, error) {
	var netClient = &http.Client{
		Timeout: time.Second * 30,
	}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	response, err := netClient.Do(req)
	return response, err
}

// PostEmptyWithHeaders wraps http.NewRequest with method POST and set request headers
func (hsc *HTTPServiceContainer) PostEmptyWithHeaders(url string, contentType string, headers map[string]string) (*http.Response, error) {
	var netClient = &http.Client{
		Timeout: time.Second * 30,
	}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)
	for headerKey, headerValue := range headers {
		req.Header.Add(headerKey, headerValue)
	}

	response, err := netClient.Do(req)
	return response, err
}
