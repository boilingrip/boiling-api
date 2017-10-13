package api

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	db, err := cleanDB()
	require.Nil(t, err)

	c, err := NewCache(db)
	require.Nil(t, err)
	require.NotNil(t, c)

	_, err = c.media.LookUp("CD")
	require.Nil(t, err)

	b := c.releaseRoles.Has("Main")
	require.True(t, b)

	keys := c.privileges.Keys()
	require.NotEmpty(t, keys)
}

func BenchmarkCacheLookUp(b *testing.B) {
	db, err := cleanDB()
	require.Nil(b, err)

	c, err := NewCache(db)
	require.Nil(b, err)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		var err error
		for pb.Next() {
			_, err = c.media.LookUp("CD")
		}
		if err != nil {
			panic(err)
		}
	})
}
