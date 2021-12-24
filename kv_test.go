package kv

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestKV(t *testing.T) {
	assert := assert.New(t)

	cache, err := newRedisTest("hello", "world", 10*time.Second)
	assert.NoError(err)
	mongo, err := newMongoTest("hello", "world")
	assert.NoError(err)

	// auto create object
	creater := func(key string) (*TestObject, error) {
		return &TestObject{
			Int: 1,
			Str: key,
		}, nil
	}

	store := NewStore(mongo, cache, creater)
	store.Del("123")
	t.Cleanup(func() {
		store.Del("123")
	})

	t.Run("Get/Set", func(t *testing.T) {
		obj, err := store.Get("123")
		assert.NoError(err)
		assert.Equal(&TestObject{Int: 1, Str: "123"}, obj)

		obj.Int = 9527
		err = store.Set("123", &obj)
		assert.NoError(err)

		obj, err = store.Get("123")
		assert.NoError(err)
		assert.Equal(&TestObject{Int: 9527, Str: "123"}, obj)
	})
}
