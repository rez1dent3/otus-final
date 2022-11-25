package fetcher

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var ErrNotSupportedContentType = errors.New("fetcher does not support content-type")

type FetchInterface interface {
	Get(context.Context, string, http.Header) ([]byte, error)
}

func NewHTTPFetcher(
	transport http.RoundTripper,
	timeout time.Duration,
	supportedContentTypes []string,
) FetchInterface {
	return &httpImpl{
		transport:             transport,
		Timeout:               timeout,
		SupportedContentTypes: supportedContentTypes,
	}
}

type httpImpl struct {
	transport http.RoundTripper
	Timeout   time.Duration

	SupportedContentTypes []string
}

func (f *httpImpl) Get(ctx context.Context, url string, header http.Header) ([]byte, error) {
	proxyRequest, err := f.prepare(ctx, url, header)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare request: %w", err)
	}

	responseBody, err := f.do(proxyRequest)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}

	return responseBody, nil
}

func (f *httpImpl) prepare(ctx context.Context, rawURL string, header http.Header) (*http.Request, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create proxy request: %w", err)
	}

	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "http://" + rawURL
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %w", err)
	}

	request.URL = parsedURL
	request.Header = header

	return request, nil
}

func (f *httpImpl) do(request *http.Request) ([]byte, error) {
	client := http.Client{
		Timeout:   f.Timeout,
		Transport: f.transport,
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to complete the request: %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	supported := false
	responseContentType := resp.Header.Get("Content-Type")
	for _, contentType := range f.SupportedContentTypes {
		if strings.Contains(responseContentType, contentType) {
			supported = true
			break
		}
	}

	if !supported {
		return nil, fmt.Errorf(
			"not supported content-type %s: %w",
			responseContentType, ErrNotSupportedContentType)
	}

	buff, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %w", err)
	}

	return buff, nil
}
