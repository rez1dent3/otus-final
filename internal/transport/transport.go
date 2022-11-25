package transport

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rez1dent3/otus-final/internal/pkg/fs"
	"github.com/rez1dent3/otus-final/internal/pkg/hsum"
	"github.com/rez1dent3/otus-final/internal/pkg/logger"
	"github.com/rez1dent3/otus-final/internal/pkg/lru"
)

var ErrServerError = errors.New("server error")

type ResponseItem struct {
	URL  string
	size uint64
}

func (i ResponseItem) Size() uint64 {
	return i.size
}

type HTTPTransport struct {
	cache lru.CacheInterface
	hash  hsum.HashInterface
	fm    fs.FileInterface
	inner http.Transport
	log   logger.LogInterface
}

func New(
	hash hsum.HashInterface,
	cache lru.CacheInterface,
	fm fs.FileInterface,
	log logger.LogInterface,
) *HTTPTransport {
	return &HTTPTransport{
		cache: cache,
		hash:  hash,
		fm:    fm,
		inner: http.Transport{},
		log:   log,
	}
}

func (t *HTTPTransport) roundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.inner.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	// redirect http to https
	if req.URL.Scheme == "http" && resp.StatusCode >= 300 && resp.StatusCode <= 399 {
		_ = resp.Body.Close()

		req.URL.Scheme = "https"

		return t.inner.RoundTrip(req)
	}

	return resp, err
}

func (t *HTTPTransport) createCache(req *http.Request) ([]byte, error) {
	now := time.Now()
	resp, err := t.roundTrip(req)
	latency := time.Since(now)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		t.log.Warning(fmt.Sprintf(
			"[%s] %s %s %s %d %s",
			now.Format("02/Jan/2006:15:04:05 -0700"),
			req.Method,
			req.URL.String(),
			resp.Status,
			latency.Microseconds(),
			req.Header.Get("User-Agent"),
		))

		return nil, ErrServerError
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if t.cache.Put(req.URL.String(), ResponseItem{
		URL:  req.URL.String(),
		size: uint64(resp.ContentLength),
	}) {
		err = t.fm.Create(t.hash.HashByString(req.URL.String()), body)
		if err != nil {
			return nil, err
		}
	}

	return body, nil
}

func (t *HTTPTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.cache.Has(req.URL.String()) {
		if body, err := t.fm.Content(t.hash.HashByString(req.URL.String())); err == nil {
			return t.response(body, err)
		}
	}

	return t.response(t.createCache(req))
}

func (t *HTTPTransport) response(body []byte, err error) (*http.Response, error) {
	if err != nil {
		return nil, err
	}

	return &http.Response{
		Header: http.Header{
			"Content-Type": []string{http.DetectContentType(body)},
		},
		Body: io.NopCloser(bytes.NewReader(body)),
	}, nil
}
