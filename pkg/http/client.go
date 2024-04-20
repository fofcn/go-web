package http

import (
	"io"
	"net/http"
)

// HTTPClient 定义了一个用于发出HTTP请求的客户端接口。
type HTTPClient interface {
	Get(url string, headers map[string]string) ([]byte, int, error)
	Post(url string, body io.Reader, headers map[string]string) ([]byte, int, error)
	Put(url string, body io.Reader, headers map[string]string) ([]byte, int, error)
	Delete(url string, headers map[string]string) ([]byte, int, error)
}

// CustomHTTPClient 实现了上面的接口并封装了HTTP请求的细节。
type CustomHTTPClient struct {
	client *http.Client
}

// NewCustomHTTPClient 返回一个新的CustomHTTPClient实例。
func NewCustomHTTPClient() *CustomHTTPClient {
	return &CustomHTTPClient{
		client: &http.Client{},
	}
}

// Get 发送一个GET请求并返回响应体和状态码。
func (c *CustomHTTPClient) Get(url string, headers map[string]string) ([]byte, int, error) {
	return c.doRequest("GET", url, nil, headers)
}

// Post 发送一个POST请求并返回响应体和状态码。
func (c *CustomHTTPClient) Post(url string, body io.Reader, headers map[string]string) ([]byte, int, error) {
	return c.doRequest("POST", url, body, headers)
}

// Put 发送一个PUT请求并返回响应体和状态码。
func (c *CustomHTTPClient) Put(url string, body io.Reader, headers map[string]string) ([]byte, int, error) {
	return c.doRequest("PUT", url, body, headers)
}

// Delete 发送一个DELETE请求并返回响应体和状态码。
func (c *CustomHTTPClient) Delete(url string, headers map[string]string) ([]byte, int, error) {
	return c.doRequest("DELETE", url, nil, headers)
}

// doRequest 是一个通用的请求处理函数，对于不同的HTTP方法封装了请求的过程。
func (c *CustomHTTPClient) doRequest(method, url string, body io.Reader, headers map[string]string) ([]byte, int, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, 0, err
	}

	if headers == nil {
		headers = make(map[string]string)
	}
	// Set request headers if any
	if headers != nil {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	return respData, resp.StatusCode, nil
}
