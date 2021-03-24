package kv

import (
	"testing"
	"time"

	"github.com/cupen/kv/mongo"
	"github.com/cupen/kv/redis"

	rds "github.com/go-redis/redis/v7"
)

func BenchmarkMongo(b *testing.B) {
	db, err := newMongoTest("kv_test", "TestObject")
	if err != nil {
		b.FailNow()
	}
	b.ResetTimer()
	b.Run("set", func(b *testing.B) {
		val := TestObject{1, "abc", 0.999, []byte{1, 2, 3}}
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			err := db.Set("abc", &val)
			if err != nil {
				b.FailNow()
			}
		}
	})

	b.Run("get", func(b *testing.B) {
		val := TestObject{1, "abc", 0.999, []byte{1, 2, 3}}
		err := db.Set("abc", &val)
		if err != nil {
			b.FailNow()
		}
		b.ResetTimer()
		b.ReportAllocs()
		val = TestObject{}
		for i := 0; i < b.N; i++ {
			err := db.Get("abc", &val)
			if err != nil {
				b.FailNow()
			}
		}
	})
}

func BenchmarkRedis(b *testing.B) {
	db, err := newRedisTest("kv_test", "TestObject", time.Minute)
	if err != nil {
		b.FailNow()
	}
	b.ResetTimer()
	b.Run("set", func(b *testing.B) {
		val := TestObject{1, "abc", 0.999, []byte{1, 2, 3}}
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			err := db.Set("abc", &val)
			if err != nil {
				b.FailNow()
			}
		}
	})

	b.Run("get", func(b *testing.B) {
		val := TestObject{1, "abc", 0.999, []byte{1, 2, 3}}
		err := db.Set("abc", &val)
		if err != nil {
			b.FailNow()
		}
		b.ResetTimer()
		b.ReportAllocs()
		val = TestObject{}
		for i := 0; i < b.N; i++ {
			err := db.Get("abc", &val)
			if err != nil {
				b.FailNow()
			}
		}
	})
}

type TestObject struct {
	Int   int
	Str   string
	Float float64
	Bytes []byte
}

func newMongoTest(database, collection string) (*mongo.Collection, error) {
	client, err := mongo.Connect("mongodb://root:root@127.0.0.1:30001/kv?authSource=admin")
	if err != nil {
		return nil, err
	}
	return mongo.New(client, database, collection)
}

func newRedisTest(keyspace, typeName string, ttl time.Duration) (*redis.Redis, error) {
	opts, err := rds.ParseURL("redis://127.0.0.1:6379/15")
	if err != nil {
		return nil, err
	}

	c := rds.NewClient(opts)
	return redis.NewRedis(c, &redis.Options{"kv_test", "TestObject", 1 * time.Minute})
}
