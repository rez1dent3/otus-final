package hsum_test

import (
	"testing"

	"github.com/rez1dent3/otus-final/pkg/hsum"
	"github.com/stretchr/testify/require"
)

func TestFnvHash_Hash(t *testing.T) {
	h := hsum.New()

	t.Run("empty", func(t *testing.T) {
		require.Equal(t, "cbf29ce484222325", h.HashByString(""))
	})

	t.Run("len5", func(t *testing.T) {
		// https://md5calc.com/hash/fnv1a64/hello
		require.Equal(t, "a430d84680aabd0b", h.HashByString("hello"))
	})

	t.Run("len100", func(t *testing.T) {
		input := "9TbCBAq4m9E0jdJCdVCr73cB2NhxBsDkoU0XgZ2lx42NvfBd4l33sVwO7sBCTWrj7Wu9RoJlepD5k8zL4rn97U49fba38zQqjYFc"

		require.Equal(t, "27f6c20de4bf8627", h.HashByString(input))
	})
}
