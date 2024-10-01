package cachex

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedisCache(t *testing.T) {
	as := assert.New(t)

	cache := NewRedisCache(RedisConfig{
		Addr: "localhost:6379",
		DB:   1,
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
		return true
	})
	as.Nil(err)

	err = cache.Close(ctx)
	as.Nil(err)
}
