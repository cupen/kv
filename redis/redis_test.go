package redis

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/qiniu/qmgo"
	"github.com/stretchr/testify/require"
)

func newRedisTest(keyspace, typeName string, ttl time.Duration) (*Redis, error) {
	opts, err := redis.ParseURL("redis://127.0.0.1:6379/15")
	if err != nil {
		return nil, err
	}

	c := redis.NewClient(opts)
	return NewRedis(c, &Options{"kv_test", "TestObject", 1 * time.Minute})
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
	assert.Equal(ErrNotFound, err)
	err = db.Set(1, &val)
	assert.NoError(err)

	val2 := TestObject1{}
	err = db.Get(2, &val2)
	assert.Equal(ErrNotFound, err)

	err = db.Get(1, &val2)
	assert.NoError(err)
	assert.Equal(val, val2)

	err = db.Del(1)
	assert.NoError(err)

	err = db.Del(2)
	assert.NoError(err)

	// t.Fatal("not implemented")
}

func TestRedis_1KDocuments(t *testing.T) {
	assert := require.New(t)

	c, err := ConnectMongo(&qmgo.Config{
		Uri: "mongodb://root:root@127.0.0.1:30001/kv?authSource=admin",
	})
	assert.NoError(err)

	db, err := NewMongoCollection(c, "kv_test", "TestObject")
	assert.NoError(err)
	t.Cleanup(func() {
		for i := 0; i < 1000; i++ {
			db.Del(i)
		}
	})
	spawn := func(i int) TestObject1 {
		return TestObject1{i, fmt.Sprintf("abc-%d", i), 0.1 * float64(i), []byte{byte(1 * i), byte(2 * i), byte(3 * i)}}
	}
	assert.NotEqual(spawn(1), spawn(2))

	for i := 0; i < 1000; i++ {
		val := spawn(i)
		err = db.Get(i, &val)
		assert.Equal(ErrNotFound, err)
		err = db.Set(i, &val)
		assert.NoError(err)

		err = db.Get(i, &val)
		assert.NoError(err)

	}

	for i := 0; i < 1000; i++ {
		val := spawn(i)
		val2 := TestObject1{}
		err = db.Get(i, &val2)
		assert.NoError(err)
		assert.Equal(val, val2)

		err = db.Del(i)
		assert.NoError(err)
		err = db.Get(i, &val2)
		assert.Equal(ErrNotFound, err)
	}
}
