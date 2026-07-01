package chainclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/fascari/token-swap-workbench/internal/config"
)

const (
	defaultTimeout    = 5 * time.Second
	maxErrorBodyBytes = 4096
)

type (
	Client struct {
		baseURL    *url.URL
		httpClient *http.Client
	}

	swapEnvelope struct {
		Swap SwapRequest `json:"Swap"`
	}
)

func New(cfg config.ChainConfig) (*Client, error) {
	baseURL, err := url.Parse(cfg.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("parsing chain base URL: %w", err)
	}

	if baseURL.Scheme == "" || baseURL.Host == "" {
		return nil, errors.New("parsing chain base URL: missing scheme or host")
	}

	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}, nil
}

func (c Client) Quote(ctx context.Context, req QuoteRequest) (QuoteResponse, error) {
	endpoint := c.endpoint("rate")
	query := endpoint.Query()
	query.Set("in", string(req.InToken))
	query.Set("out", string(req.OutToken))
	query.Set("amount", strconv.FormatFloat(req.Amount, 'f', -1, 64))
	endpoint.RawQuery = query.Encode()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return QuoteResponse{}, fmt.Errorf("building chain quote request: %w", err)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return QuoteResponse{}, fmt.Errorf("requesting chain quote: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if err := ensureSuccess(resp); err != nil {
		return QuoteResponse{}, fmt.Errorf("requesting chain quote: %w", err)
	}

	var quote QuoteResponse
	if err := json.NewDecoder(resp.Body).Decode(&quote); err != nil {
		return QuoteResponse{}, fmt.Errorf("decoding chain quote response: %w", err)
	}

	return quote, nil
}

func (c Client) SubmitSwap(ctx context.Context, req SwapRequest) error {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(swapEnvelope{Swap: req}); err != nil {
		return fmt.Errorf("encoding chain swap transaction: %w", err)
	}

	endpoint := c.endpoint("transaction")
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), &body)
	if err != nil {
		return fmt.Errorf("building chain swap request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("submitting chain swap: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if err := ensureSuccess(resp); err != nil {
		return fmt.Errorf("submitting chain swap: %w", err)
	}

	return nil
}

func (c Client) Blocks(ctx context.Context, n int) ([]Block, error) {
	endpoint := c.endpoint("blocks")
	query := endpoint.Query()
	query.Set("n", strconv.Itoa(n))
	endpoint.RawQuery = query.Encode()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("building chain blocks request: %w", err)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("requesting chain blocks: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if err := ensureSuccess(resp); err != nil {
		return nil, fmt.Errorf("requesting chain blocks: %w", err)
	}

	var blocks []Block
	if err := json.NewDecoder(resp.Body).Decode(&blocks); err != nil {
		return nil, fmt.Errorf("decoding chain blocks response: %w", err)
	}

	return blocks, nil
}

func (c Client) Status(ctx context.Context) error {
	if _, err := c.Blocks(ctx, 1); err != nil {
		return fmt.Errorf("checking chain status: %w", err)
	}

	return nil
}

func (c Client) endpoint(route string) url.URL {
	endpoint := *c.baseURL
	endpoint.Path = strings.TrimRight(endpoint.Path, "/") + "/" + strings.TrimLeft(route, "/")
	endpoint.RawQuery = ""
	return endpoint
}

func ensureSuccess(resp *http.Response) error {
	if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices {
		return nil
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxErrorBodyBytes))
	if err != nil {
		return fmt.Errorf("chain returned %s and the error body could not be read: %w", resp.Status, err)
	}

	detail := strings.TrimSpace(string(body))
	if detail == "" {
		return fmt.Errorf("chain returned %s", resp.Status)
	}

	return fmt.Errorf("chain returned %s: %s", resp.Status, detail)
}
