package bytesize_test

import (
	"testing"

	"github.com/rez1dent3/otus-final/pkg/bytesize"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	t.Run("length of number 1", func(t *testing.T) {
		require.Equal(t, uint64(1), bytesize.Parse("1B"))
		require.Equal(t, uint64(1*1024), bytesize.Parse("1K"))
		require.Equal(t, uint64(1*1024*1024), bytesize.Parse("1M"))
		require.Equal(t, uint64(1*1024*1024*1024), bytesize.Parse("1G"))
	})

	t.Run("Length of number 2+", func(t *testing.T) {
		require.Equal(t, uint64(128), bytesize.Parse("128B"))
		require.Equal(t, uint64(128*1024), bytesize.Parse("128K"))
		require.Equal(t, uint64(128*1024*1024), bytesize.Parse("128M"))
		require.Equal(t, uint64(128*1024*1024*1024), bytesize.Parse("128G"))
	})

	t.Run("errors", func(t *testing.T) {
		require.Equal(t, uint64(0), bytesize.Parse("."))
		require.Equal(t, uint64(0), bytesize.Parse("1KB"))
		require.Equal(t, uint64(0), bytesize.Parse("hello"))
		require.Equal(t, uint64(0), bytesize.Parse("1.5K"))
	})
}
