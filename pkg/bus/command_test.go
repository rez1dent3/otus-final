package bus_test

import (
	"testing"

	"github.com/rez1dent3/otus-final/pkg/bus"
	"github.com/stretchr/testify/require"
)

func TestBus(t *testing.T) {
	t.Run("fire.unsubscribe", func(t *testing.T) {
		b := bus.NewSyncBus()
		b.Fire("hello", "world")
		require.True(t, true)
	})

	t.Run("fire.single", func(t *testing.T) {
		c := 0
		b := bus.NewSyncBus()

		b.Subscribe("single", func(a any) {
			if val, ok := a.(int); ok {
				c += val
			}
		})

		require.Equal(t, 0, c)

		b.Fire("single", 5)
		require.Equal(t, 5, c)

		b.Fire("single", 3)
		require.Equal(t, 8, c)
	})

	t.Run("fire.multi", func(t *testing.T) {
		c1 := 0
		c2 := 1
		b := bus.NewSyncBus()

		// add
		b.Subscribe("multi", func(a any) {
			if val, ok := a.(int); ok {
				c1 += val
			}
		})

		// pow
		b.Subscribe("multi", func(a any) {
			if val, ok := a.(int); ok {
				c2 *= val
			}
		})

		require.Equal(t, 0, c1)
		require.Equal(t, 1, c2)

		b.Fire("multi", 5)
		require.Equal(t, 5, c1)
		require.Equal(t, 5, c2)

		b.Fire("multi", 5)
		require.Equal(t, 10, c1)
		require.Equal(t, 25, c2)
	})
}
