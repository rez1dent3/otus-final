package transformer_test

import (
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/rez1dent3/otus-final/internal/pkg/transformer"
	"github.com/stretchr/testify/require"
)

type text struct{}

func (t *text) FillCenter(source []byte, _, _ int) ([]byte, error) {
	return source, nil
}

func (t *text) IsSupported(source []byte) bool {
	return strings.Contains(http.DetectContentType(source), "text/plain")
}

func TestStack_IsSupported(t *testing.T) {
	transform := transformer.NewStackBy(transformer.NewJpeg(), transformer.NewPng(), &text{})

	testCases := []struct {
		name      string
		path      string
		supported bool
	}{
		{"check jpeg", "../../../resources/images/_gopher_original_1024x504.jpg", true},
		{"check png", "../../../resources/images/_gopher_original_1024x504.png", true},
		{"check go", "./fill_jpeg.go", true},
		{"check /dev/null", "/dev/null", true},
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

func TestStack_FillCenter(t *testing.T) {
	transform := transformer.NewStackBy(transformer.NewJpeg(), &text{})

	testCases := []struct {
		name string
		path string
	}{
		{"check /bin/sh", "/bin/sh"},
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

			_, err = transform.FillCenter(readAll, 1, 1)
			require.ErrorIs(t, transformer.ErrFileNotSupported, err)
		})
	}
}
