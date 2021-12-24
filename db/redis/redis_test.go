package redis

import (
	"testing"
	"time"

	"github.com/cupen/kv/errors"
	redis "github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/require"
)

func newRedisTest(keyspace, typeName string, ttl time.Duration) (*Redis, error) {
	opts, err := redis.ParseURL("redis://127.0.0.1:6379/15")
	if err != nil {
		return nil, err
	}
	c := redis.NewClient(opts)
	return New(c, &Options{"kv_test", "TestObject", 1 * time.Minute})
}

type TestObject1 struct {
	Int   int
	Str   string
	Float float64
	Bytes []byte
}

func TestRedis(t *testing.T) {
	assert := require.New(t)

	var db, err = newRedisTest("keyspace", "TestObject1", 1*time.Minute)
	assert.NoError(err)

	t.Cleanup(func() {
		db.Del(1)
		db.Del(2)
	})

	val := TestObject1{1, "abc", 0.999, []byte{1, 2, 3}}
	err = db.Get(1, &val)
	assert.Equal(errors.ErrNotFound, err)
	err = db.Set(1, &val)
	assert.NoError(err)

	val2 := TestObject1{}
	err = db.Get(2, &val2)
	assert.Equal(errors.ErrNotFound, err)

	err = db.Get(1, &val2)
	assert.NoError(err)
	assert.Equal(val, val2)

	err = db.Del(1)
	assert.NoError(err)

	err = db.Del(2)
	assert.NoError(err)

	// t.Fatal("not implemented")
}
