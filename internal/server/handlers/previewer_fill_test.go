package handlers_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/rez1dent3/otus-final/internal/server/handlers"
	"github.com/stretchr/testify/require"
)

var img = "raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/_gopher_original_1024x504.jpg" //nolint:lll

func TestPreviewHandler_ParseURL(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ph := handlers.PreviewHandler{}
		originalURL, width, height, err := ph.ParseURL(&http.Request{URL: &url.URL{Path: "/fill/100/100/" + img}})
		require.NoError(t, err)
		require.Equal(t, img, originalURL)
		require.Equal(t, 100, width)
		require.Equal(t, 100, height)
	})
}
