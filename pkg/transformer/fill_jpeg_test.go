//nolint:dupl
package transformer_test

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/rez1dent3/otus-final/pkg/hsum"
	"github.com/rez1dent3/otus-final/pkg/transformer"
	"github.com/stretchr/testify/require"
)

func TestJpegImage_IsSupported(t *testing.T) {
	transform := transformer.NewJpeg()

	testCases := []struct {
		name      string
		path      string
		supported bool
	}{
		{"check jpeg", "../../resources/images/_gopher_original_1024x504.jpg", true},
		{"check png", "../../resources/images/_gopher_original_1024x504.png", false},
		{"check go", "./fill_jpeg.go", false},
		{"check /dev/null", "/dev/null", false},
		{"check /bin/sh", "/bin/sh", false},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			file, err := os.Open(testCase.path)
			require.NoError(t, err)
			defer func() {
				_ = file.Close()
			}()

			readAll, err := io.ReadAll(file)
			require.NoError(t, err)

			require.Equal(t, testCase.supported, transform.IsSupported(readAll))
		})
	}
}

func TestJpegImage_FillCenter(t *testing.T) {
	// prepare
	transform := transformer.NewJpeg()
	hsm := hsum.New()

	file, err := os.Open("../../resources/images/_gopher_original_1024x504.jpg")
	require.NoError(t, err)
	require.NotNil(t, file)

	defer func() {
		_ = file.Close()
	}()
	readAll, err := io.ReadAll(file)
	require.NoError(t, err)

	testCases := []struct {
		width, height int
		expected      string
	}{
		{50, 50, "75f6e244e52b5aca"},
		{200, 700, "2bad63c71b9e5903"},
		{500, 500, "b01ec26ccedb2004"},
		{1024, 252, "3a2e354ae5f7ca23"},
		{1025, 600, "2afa38a81fa0b377"},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(fmt.Sprintf("fill %dx%d", testCase.width, testCase.height), func(t *testing.T) {
			dst, err := transform.FillCenter(readAll, testCase.width, testCase.height)

			require.NoError(t, err)
			require.Equal(t, testCase.expected, hsm.Hash(dst))
		})
	}
}
