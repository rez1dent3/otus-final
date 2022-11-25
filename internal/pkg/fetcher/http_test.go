package fetcher_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/rez1dent3/otus-final/internal/pkg/fetcher"
	"github.com/rez1dent3/otus-final/internal/pkg/hsum"
	"github.com/stretchr/testify/require"
)

func checkAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "test" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func fileServer() *httptest.Server {
	return httptest.NewServer(checkAuth(http.FileServer(http.Dir("../../../resources/images"))))
}

func TestHttpImpl_Get(t *testing.T) {
	t.Run("httpauth", func(t *testing.T) {
		testCases := []struct {
			headers        http.Header
			expectedStatus int
		}{
			{nil, http.StatusUnauthorized},
			{http.Header{"Authorization": []string{"test"}}, http.StatusOK},
		}

		server := fileServer()
		defer server.Close()

		for _, testCase := range testCases {
			req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL, nil)
			req.Header = testCase.headers
			require.NoError(t, err)

			client := http.Client{}
			response, err := client.Do(req)
			_ = response.Body.Close()

			require.NoError(t, err)

			require.Equal(t, testCase.expectedStatus, response.StatusCode)
		}
	})

	t.Run("images not supported", func(t *testing.T) {
		testCases := []struct {
			name string
		}{
			{"_gopher_original_1024x504.webp"},
		}

		fetch := fetcher.NewHTTPFetcher(
			&http.Transport{},
			50*time.Millisecond,
			[]string{"image/jpeg", "image/png"},
		)

		server := fileServer()
		defer server.Close()

		for _, testCase := range testCases {
			rawURL, err := url.JoinPath(server.URL, testCase.name)
			require.NoError(t, err)
			require.Equal(t, server.URL+"/"+testCase.name, rawURL)

			_, err = fetch.Get(context.Background(), rawURL, http.Header{
				"Authorization": []string{"test"},
			})

			require.ErrorIs(t, err, fetcher.ErrNotSupportedContentType)
		}
	})

	t.Run("images supported", func(t *testing.T) {
		testCases := []struct {
			name     string
			expected string
		}{
			{"_gopher_original_1024x504.jpg", "f14b75d39d77d92b"},
			{"_gopher_original_1024x504.png", "52100d3fa2c9c200"},
		}

		fetch := fetcher.NewHTTPFetcher(
			&http.Transport{},
			50*time.Millisecond,
			[]string{"image/jpeg", "image/png"},
		)

		server := fileServer()
		defer server.Close()

		for _, testCase := range testCases {
			rawURL, err := url.JoinPath(server.URL, testCase.name)
			require.NoError(t, err)
			require.Equal(t, server.URL+"/"+testCase.name, rawURL)

			bytes, err := fetch.Get(context.Background(), rawURL, http.Header{
				"Authorization": []string{"test"},
			})

			require.NoError(t, err)
			require.Equal(t, testCase.expected, hsum.New().Hash(bytes))
		}
	})
}
