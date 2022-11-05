package tests

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/disintegration/imaging"
	"github.com/stretchr/testify/require"
)

func doRequest(rawURL string, width, height int, header http.Header) (*http.Response, error) {
	req, _ := http.NewRequestWithContext(context.Background(), "GET", "http://imgproxy:8000", nil)
	req.URL.Path = fmt.Sprintf("/fill/%d/%d/%s", width, height, rawURL)
	req.Header = header

	return http.DefaultClient.Do(req)
}

func TestCheckErrors(t *testing.T) {
	testCases := []struct {
		url           string
		width, height int
	}{
		// remote server does not exist
		{"domain-not-exists/gopher.png", 640, 480},

		// the remote server exists, but the image was not found
		{"nginx/4xx", 640, 480},

		// check not supported formats
		{"nginx/_gopher_original_1024x504.webp", 640, 480},
		{"nginx/text", 640, 480},

		// the remote server returned an error
		{"nginx/5xx", 640, 480},
	}

	for _, testCase := range testCases {
		resp, _ := doRequest(testCase.url, testCase.width, testCase.height, nil)
		require.NotNil(t, resp)

		require.Equal(t, http.StatusBadGateway, resp.StatusCode)
		require.NoError(t, resp.Body.Close())
	}
}

func TestCheckSupportImages(t *testing.T) {
	testCases := []struct {
		url           string
		width, height int
	}{
		// check support formats
		{"nginx/_gopher_original_1024x504.jpg", 640, 480},
		{"nginx/_gopher_original_1024x504.png", 640, 480},

		// limit values
		{"nginx/_gopher_original_1024x504.jpg", 1, 1},
		{"nginx/_gopher_original_1024x504.jpg", 1024, 504},
		{"nginx/_gopher_original_1024x504.jpg", 4000, 2000},
	}

	for _, testCase := range testCases {
		resp, _ := doRequest(testCase.url, testCase.width, testCase.height, nil)
		require.NotNil(t, resp)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		image, err := imaging.Decode(resp.Body)
		require.NoError(t, err)

		require.Equal(t, testCase.width, image.Bounds().Dx())
		require.Equal(t, testCase.height, image.Bounds().Dy())

		require.NoError(t, resp.Body.Close())
	}
}

func TestCheckContentTypeHeader(t *testing.T) {
	testCases := []struct {
		url           string
		width, height int
		header        string
	}{
		// check support formats
		{"nginx/_gopher_original_1024x504.jpg", 640, 480, "image/jpeg"},
		{"nginx/_gopher_original_1024x504.png", 640, 480, "image/png"},
	}

	for _, testCase := range testCases {
		resp, _ := doRequest(testCase.url, testCase.width, testCase.height, nil)
		require.NotNil(t, resp)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.Equal(t, testCase.header, resp.Header.Get("Content-Type"))
		require.NoError(t, resp.Body.Close())
	}
}

func TestCheckFromCachePreview(t *testing.T) {
	// first request without cache
	resp, _ := doRequest("nginx/limited/_gopher_original_1024x504.jpg", 640, 480, nil)
	require.NotNil(t, resp)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.NoError(t, resp.Body.Close())

	// second request without cache (new parameters for resizing). check nginx rate limit
	resp, _ = doRequest("nginx/limited/_gopher_original_1024x504.jpg", 1024, 768, nil)
	require.NotNil(t, resp)
	require.Equal(t, http.StatusBadGateway, resp.StatusCode)
	require.NoError(t, resp.Body.Close())

	// third request with cache
	resp, _ = doRequest("nginx/limited/_gopher_original_1024x504.jpg", 640, 480, nil)
	require.NotNil(t, resp)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.NoError(t, resp.Body.Close())
}

func TestCheckSendingHeaders(t *testing.T) {
	// first request without auth
	resp, _ := doRequest("nginx/auth/_gopher_original_1024x504.jpg", 640, 480, nil)
	require.NotNil(t, resp)
	require.Equal(t, http.StatusBadGateway, resp.StatusCode)
	require.NoError(t, resp.Body.Close())

	// second request with auth
	resp, _ = doRequest("nginx/auth/_gopher_original_1024x504.jpg", 640, 480, http.Header{
		"Authorization": []string{"Basic dXNlcjp1c2Vy"},
	})
	require.NotNil(t, resp)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.NoError(t, resp.Body.Close())
}
