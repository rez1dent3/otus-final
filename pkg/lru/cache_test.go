package lru_test

import (
	"testing"

	"github.com/rez1dent3/otus-final/pkg/bus"
	"github.com/rez1dent3/otus-final/pkg/lru"
	"github.com/stretchr/testify/require"
)

type val struct {
	size uint64
}

func (v val) Size() uint64 {
	return v.size
}

func TestLru_Limits(t *testing.T) {
	t.Run("limit<size", func(t *testing.T) {
		c := lru.New(0, bus.NewSyncBus())
		require.False(t, c.Put("hello", val{5}))

		res, ok := c.Get("hello")
		require.Nil(t, res)
		require.False(t, ok)

		require.Equal(t, uint64(0), c.Size())
	})

	t.Run("limit>size", func(t *testing.T) {
		c := lru.New(5, bus.NewSyncBus())
		require.True(t, c.Put("hello", val{4}))

		res, ok := c.Get("hello")
		require.Equal(t, uint64(4), res.(val).size)
		require.True(t, ok)

		require.Equal(t, uint64(4), c.Size())
	})

	t.Run("limit=size", func(t *testing.T) {
		c := lru.New(5, bus.NewSyncBus())
		require.True(t, c.Put("hello", val{3}))

		res, ok := c.Get("hello")
		require.Equal(t, uint64(3), res.(val).size)
		require.True(t, ok)

		require.True(t, c.Put("hello", val{5}))

		res, ok = c.Get("hello")
		require.Equal(t, uint64(5), res.(val).size)
		require.True(t, ok)

		require.Equal(t, uint64(5), c.Size())
	})

	t.Run("limit=0,size=0", func(t *testing.T) {
		c := lru.New(0, bus.NewSyncBus())
		require.True(t, c.Put("hello", val{0}))
		require.True(t, c.Put("world", val{0}))

		require.True(t, c.Has("hello"))
		require.True(t, c.Has("world"))

		require.Equal(t, uint64(0), c.Size())

		res, ok := c.Get("hello")
		require.Equal(t, uint64(0), res.(val).size)
		require.True(t, ok)

		res, ok = c.Get("world")
		require.Equal(t, uint64(0), res.(val).size)
		require.True(t, ok)

		require.Equal(t, uint64(0), c.Size())
	})

	t.Run("limit=5,stream=3,2,2,2", func(t *testing.T) {
		c := lru.New(5, bus.NewSyncBus())
		require.True(t, c.Put("a", val{3}))
		require.True(t, c.Put("b", val{2}))
		require.True(t, c.Put("c", val{2}))
		require.True(t, c.Put("d", val{2}))

		require.False(t, c.Has("a"))
		require.False(t, c.Has("b"))
		require.True(t, c.Has("c"))
		require.True(t, c.Has("d"))

		require.Equal(t, uint64(4), c.Size())
	})

	t.Run("update item", func(t *testing.T) {
		c := lru.New(5, bus.NewSyncBus())
		require.True(t, c.Put("hello", val{5}))

		res, ok := c.Get("hello")
		require.Equal(t, uint64(5), res.(val).size)
		require.True(t, ok)

		require.True(t, c.Put("hello", val{4}))
		require.True(t, c.Put("world", val{1}))

		res, ok = c.Get("hello")
		require.Equal(t, uint64(4), res.(val).size)
		require.True(t, ok)

		res, ok = c.Get("world")
		require.Equal(t, uint64(1), res.(val).size)
		require.True(t, ok)
	})
}

func TestLru_Evict(t *testing.T) {
	t.Run("evict", func(t *testing.T) {
		c := lru.New(5, bus.NewSyncBus())
		require.True(t, c.Put("a", val{1}))
		require.True(t, c.Put("b", val{1}))
		require.True(t, c.Put("c", val{1}))
		require.True(t, c.Put("d", val{1}))
		require.True(t, c.Put("e", val{1}))
		require.True(t, c.Put("f", val{1}))

		// [f,e,d,c,b]
		require.False(t, c.Has("a"))
		require.True(t, c.Has("b"))
		require.True(t, c.Has("c"))
		require.True(t, c.Has("d"))
		require.True(t, c.Has("e"))
		require.True(t, c.Has("f"))

		res, ok := c.Get("c")
		require.Equal(t, uint64(1), res.(val).size)
		require.True(t, ok)

		// [c,f,e,d,b]
		require.True(t, c.Has("c"))
		require.True(t, c.Has("b"))
		require.True(t, c.Has("d"))
		require.True(t, c.Has("e"))
		require.True(t, c.Has("f"))

		res, ok = c.Get("f")
		require.Equal(t, uint64(1), res.(val).size)
		require.True(t, ok)

		// [f,c,e,d,b]
		require.True(t, c.Has("f"))
		require.True(t, c.Has("c"))
		require.True(t, c.Has("b"))
		require.True(t, c.Has("d"))
		require.True(t, c.Has("e"))

		require.True(t, c.Put("g", val{1}))

		// [g,f,c,e,d]
		require.True(t, c.Has("g"))
		require.True(t, c.Has("f"))
		require.True(t, c.Has("c"))
		require.True(t, c.Has("e"))
		require.True(t, c.Has("d"))

		// removed
		require.False(t, c.Has("a"))
		require.False(t, c.Has("b"))

		c.Purge()

		require.False(t, c.Has("g"))
		require.False(t, c.Has("f"))
		require.False(t, c.Has("c"))
		require.False(t, c.Has("e"))
		require.False(t, c.Has("d"))
	})
}
