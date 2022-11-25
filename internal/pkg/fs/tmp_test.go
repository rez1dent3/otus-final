package fs_test

import (
	"os"
	"testing"

	"github.com/rez1dent3/otus-final/internal/pkg/fs"
	"github.com/stretchr/testify/require"
)

func TestFm(t *testing.T) {
	fm := fs.New(os.TempDir(), "test")

	t.Run("simple", func(t *testing.T) {
		require.NoError(t, fm.Create("hello", []byte("hello world")))

		cnt, err := fm.Content("hello")
		require.NoError(t, err)
		require.Equal(t, []byte("hello world"), cnt)

		require.NoError(t, fm.Delete("hello"))

		cnt, err = fm.Content("hello")
		require.ErrorIs(t, fs.ErrOpenFile, err)
		require.Nil(t, nil, cnt)
	})

	t.Run("replace data", func(t *testing.T) {
		require.NoError(t, fm.Create("hello", []byte("hello world")))

		cnt, err := fm.Content("hello")
		require.NoError(t, err)
		require.Equal(t, []byte("hello world"), cnt)

		require.NoError(t, fm.Create("hello", []byte("hello")))

		cnt, err = fm.Content("hello")
		require.NoError(t, err)
		require.Equal(t, []byte("hello"), cnt)
	})

	t.Run("failed.create", func(t *testing.T) {
		fm := fs.New("/dev", "failed")
		require.ErrorIs(t, fs.ErrCreateFile, fm.Create("hello", []byte("")))
	})

	t.Run("failed.delete", func(t *testing.T) {
		require.ErrorIs(t, fs.ErrDeleteFile, fm.Delete("failed-del"))
	})
}
