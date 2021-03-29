package kv

import (
	"testing"
	"time"
)

func BenchmarkMongoCollection(b *testing.B) {
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
