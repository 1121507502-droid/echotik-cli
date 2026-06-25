package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	BaseURL    string
	Username   string
	Password   string
	HTTPClient *http.Client
}

type Request struct {
	Method string
	Path   string
	Params map[string]string
	Body   any
}

type Response struct {
	StatusCode int
	Header     http.Header
	Raw        []byte
	JSON       any
}

func New(baseURL, username, password string) *Client {
	return &Client{
		BaseURL:  strings.TrimRight(baseURL, "/"),
		Username: username,
		Password: password,
		HTTPClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (c *Client) Do(ctx context.Context, r Request) (*Response, error) {
	endpoint, err := c.buildURL(r.Path, r.Params)
	if err != nil {
		return nil, err
	}

	var body io.Reader
	if r.Body != nil {
		b, err := json.Marshal(r.Body)
		if err != nil {
			return nil, fmt.Errorf("encode request body: %w", err)
		}
		body = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, strings.ToUpper(r.Method), endpoint, body)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.Username, c.Password)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, 32*1024*1024))
	if err != nil {
		return nil, err
	}
	out := &Response{StatusCode: resp.StatusCode, Header: resp.Header, Raw: raw}
	if len(raw) > 0 && strings.Contains(resp.Header.Get("Content-Type"), "json") {
		var v any
		if err := json.Unmarshal(raw, &v); err == nil {
			out.JSON = v
		}
	}
	if resp.StatusCode >= 400 {
		return out, &HTTPError{StatusCode: resp.StatusCode, Body: string(raw)}
	}
	return out, nil
}

func (c *Client) buildURL(path string, params map[string]string) (string, error) {
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path, nil
	}
	u, err := url.Parse(c.BaseURL + "/" + strings.TrimLeft(path, "/"))
	if err != nil {
		return "", err
	}
	q := u.Query()
	for k, v := range params {
		if v != "" {
			q.Set(k, v)
		}
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}

type HTTPError struct {
	StatusCode int
	Body       string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, strings.TrimSpace(e.Body))
}
