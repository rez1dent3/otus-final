package transport_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/rez1dent3/otus-final/internal/pkg/bus"
	"github.com/rez1dent3/otus-final/internal/pkg/bytesize"
	"github.com/rez1dent3/otus-final/internal/pkg/fs"
	"github.com/rez1dent3/otus-final/internal/pkg/hsum"
	"github.com/rez1dent3/otus-final/internal/pkg/logger"
	"github.com/rez1dent3/otus-final/internal/pkg/lru"
	"github.com/rez1dent3/otus-final/internal/transport"
	"github.com/stretchr/testify/require"
)

func fileServer() *httptest.Server {
	return httptest.NewServer(http.FileServer(http.Dir("../../resources/images")))
}

func newCache(hash hsum.HashInterface, fm fs.FileInterface) lru.CacheInterface {
	commandBus := bus.NewSyncBus()
	cache := lru.New(bytesize.Parse("65K"), commandBus)

	commandBus.Subscribe(lru.EventEvict, func(arg any) {
		if data, ok := arg.(transport.ResponseItem); ok {
			_ = fm.Delete(hash.HashByString(data.URL))
		}
	})

	return cache
}

func newTransport(
	hash hsum.HashInterface,
	fm fs.FileInterface,
	cache lru.CacheInterface,
) *transport.HTTPTransport {
	return transport.New(hash, cache, fm, logger.New("off", nil))
}

func TestHTTPTransport_RoundTrip(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		hash := hsum.New()
		fm := fs.New(os.TempDir(), "transport-test")
		cache := newCache(hash, fm)

		defer cache.Purge()

		testCases := []struct {
			name      string
			expected  string
			fromCache bool
			header    http.Header
		}{
			// size < limit
			{"_gopher_original_1024x504.jpg", "f14b75d39d77d92b", false, http.Header{
				"Content-Type": []string{"image/jpeg"},
			}},
			{"_gopher_original_1024x504.jpg", "f14b75d39d77d92b", true, http.Header{
				"Content-Type": []string{"image/jpeg"},
			}},

			// size > limit
			{"_gopher_original_1024x504.png", "52100d3fa2c9c200", false, http.Header{
				"Content-Type": []string{"image/png"},
			}},
			{"_gopher_original_1024x504.png", "52100d3fa2c9c200", false, http.Header{
				"Content-Type": []string{"image/png"},
			}},
		}

		client := http.Client{Transport: newTransport(hash, fm, cache), Timeout: time.Second}

		server := fileServer()
		defer server.Close()

		for _, testCase := range testCases {
			rawURL, err := url.JoinPath(server.URL, testCase.name)
			require.NoError(t, err)
			require.Equal(t, server.URL+"/"+testCase.name, rawURL)

			req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, rawURL, nil)
			req.Header = http.Header{}
			require.NoError(t, err)

			fromCache := cache.Has(rawURL)
			response, err := client.Do(req)
			require.NoError(t, err)

			readAll, err := io.ReadAll(response.Body)
			require.NoError(t, err)

			err = response.Body.Close()
			require.NoError(t, err)

			require.NoError(t, err)
			require.Equal(t, testCase.expected, hsum.New().Hash(readAll))
			require.Equal(t, testCase.header, response.Header)
			require.Equal(t, testCase.fromCache, fromCache)
		}
	})
}
