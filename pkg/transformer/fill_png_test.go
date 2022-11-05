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

func TestPngImage_IsSupported(t *testing.T) {
	transformPng := transformer.NewPng()

	testCases := []struct {
		name      string
		path      string
		supported bool
	}{
		{"check jpeg", "../../resources/images/_gopher_original_1024x504.jpg", false},
		{"check png", "../../resources/images/_gopher_original_1024x504.png", true},
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

			require.Equal(t, testCase.supported, transformPng.IsSupported(readAll))
		})
	}
}

func TestPngImage_FillCenter(t *testing.T) {
	// prepare
	transformPng := transformer.NewPng()
	hsm := hsum.New()

	file, err := os.Open("../../resources/images/_gopher_original_1024x504.png")
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
		{50, 50, "20deb5863fdb018c"},
		{200, 700, "102caa63032b64f2"},
		{500, 500, "33b6db1de3192b1a"},
		{1024, 252, "31715390d154954d"},
		{1025, 600, "06ac5873c4fac3e6"},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(fmt.Sprintf("fill %dx%d", testCase.width, testCase.height), func(t *testing.T) {
			dst, err := transformPng.FillCenter(readAll, testCase.width, testCase.height)

			require.NoError(t, err)
			require.Equal(t, testCase.expected, hsm.Hash(dst))
		})
	}
}
