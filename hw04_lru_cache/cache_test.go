package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	//nolint:depguard // Применение 'require' необходимо для тестирования.
	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(2)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		c.Clear()

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("displacement oldest logic", func(t *testing.T) {
		c := NewCache(3)

		wasInCache := c.Set("one", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("two", 200)
		require.False(t, wasInCache)

		wasInCache = c.Set("three", 300)
		require.False(t, wasInCache)

		val, ok := c.Get("one")
		require.True(t, ok)
		require.Equal(t, 100, val)

		wasInCache = c.Set("four", 400)
		require.False(t, wasInCache)

		_, ok = c.Get("two")
		require.False(t, ok)

		_, ok = c.Get("one")
		require.True(t, ok)
	})

	t.Run("displacement not used logic", func(t *testing.T) {
		c := NewCache(3)

		wasInCache := c.Set("one", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("two", 200)
		require.False(t, wasInCache)

		wasInCache = c.Set("three", 300)
		require.False(t, wasInCache)

		wasInCache = c.Set("four", 400)
		require.False(t, wasInCache)

		wasInCache = c.Set("five", 500)
		require.False(t, wasInCache)

		_, ok := c.Get("one")
		require.False(t, ok)

		_, ok = c.Get("five")
		require.True(t, ok)
	})
}

func TestCacheMultithreading(t *testing.T) {
	t.Skip() // Remove me if task with asterisk completed.

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
