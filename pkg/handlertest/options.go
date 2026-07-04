package handlertest

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
)

type (
	OptionFunc func(*requestOptions)

	requestOptions struct {
		body    io.Reader
		query   url.Values
		headers http.Header
		ctx     context.Context
	}
)

func WithBody(body string) OptionFunc {
	return func(o *requestOptions) {
		o.body = strings.NewReader(body)
	}
}

func WithQuery(key string, value string) OptionFunc {
	return func(o *requestOptions) {
		o.query.Set(key, value)
	}
}

func WithHeader(key string, value string) OptionFunc {
	return func(o *requestOptions) {
		o.headers.Set(key, value)
	}
}

func WithContext(ctx context.Context) OptionFunc {
	return func(o *requestOptions) {
		o.ctx = ctx
	}
}

func (s *Suite) buildRequest(method string, target string, opts ...OptionFunc) *http.Request {
	s.T().Helper()

	options := requestOptions{
		query:   url.Values{},
		headers: http.Header{},
		ctx:     s.T().Context(),
	}

	for _, opt := range opts {
		opt(&options)
	}

	if encoded := options.query.Encode(); encoded != "" {
		target += "?" + encoded
	}

	request := httptest.NewRequestWithContext(options.ctx, method, target, options.body)
	for key, values := range options.headers {
		request.Header[key] = values
	}

	return request
}
