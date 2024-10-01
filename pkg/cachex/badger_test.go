package cachex

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBadgerCache(t *testing.T) {
	as := assert.New(t)

	cache := NewBadgerCache(BadgerConfig{
		Path: "./tmp/badger",
	})

	ctx := context.Background()
	err := cache.Set(ctx, "tt", "foo", "bar")
	as.Nil(err)

	val, exists, err := cache.Get(ctx, "tt", "foo")
	as.Nil(err)
	as.True(exists)
	as.Equal("bar", val)

	err = cache.Delete(ctx, "tt", "foo")
	as.Nil(err)

	val, exists, err = cache.Get(ctx, "tt", "foo")
	as.Nil(err)
	as.False(exists)
	as.Equal("", val)

	tmap := make(map[string]bool)
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("foo%d", i)
		err = cache.Set(ctx, "tt", key, "bar")
		as.Nil(err)
		tmap[key] = true

		err = cache.Set(ctx, "ff", key, "bar")
		as.Nil(err)
	}

	err = cache.Iterator(ctx, "tt", func(ctx context.Context, key, value string) bool {
		as.True(tmap[key])
		as.Equal("bar", value)
		t.Log(key, value)
		return true
	})
	as.Nil(err)

	err = cache.Close(ctx)
	as.Nil(err)
}
