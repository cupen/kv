package mongo

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/cupen/kv/errors"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

type TestObject struct {
	Int   int
	Str   string
	Float float64
	Bytes []byte
}

func newMongoTest(database, collection string) (*Collection, error) {
	client, err := Connect("mongodb://root:root@127.0.0.1:30001/kv?authSource=admin")
	if err != nil {
		return nil, err
	}
	return New(client, database, collection)
}

func TestMongoCollection(t *testing.T) {
	assert := require.New(t)

	db, err := newMongoTest("kv_test", "TestObject")
	assert.NoError(err)
	t.Cleanup(func() {
		db.Del(1)
		db.Del(2)
	})
	ErrNotFound := errors.ErrNotFound

	val := TestObject{1, "abc", 0.999, []byte{1, 2, 3}}
	err = db.Get(1, &val)
	assert.Equal(ErrNotFound, err)
	err = db.Set(1, &val)
	assert.NoError(err)

	val2 := TestObject{}
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

func TestMongoCollection_1KDocuments(t *testing.T) {
	assert := require.New(t)
	ErrNotFound := errors.ErrNotFound

	db, err := newMongoTest("kv_test", "TestObject_1kdocs")
	assert.NoError(err)
	t.Cleanup(func() {
		db.coll.DropCollection(context.TODO())
	})
	spawn := func(i int) TestObject {
		return TestObject{i, fmt.Sprintf("abc-%d", i), 0.1 * float64(i), []byte{byte(1 * i), byte(2 * i), byte(3 * i)}}
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
		val2 := TestObject{}
		err = db.Get(i, &val2)
		assert.NoError(err)
		assert.Equal(val, val2)

		err = db.Del(i)
		assert.NoError(err)
		err = db.Get(i, &val2)
		assert.Equal(ErrNotFound, err)
	}
}

func BenchmarkMongoCollection(b *testing.B) {
	newObject := func(i int) *TestObject {
		return &TestObject{
			Int:   i,
			Str:   strings.Repeat(strconv.Itoa(i), 10),
			Float: float64(i),
		}
	}
	newObjects := func(start, end int) []*TestObject {
		size := end - start
		if size <= 0 || size > 1000 {
			panic(fmt.Errorf("invalid range: start:%d end:%d", start, end))
		}
		rs := make([]*TestObject, size)
		for i := 0; i < size; i++ {
			tag := start + i
			rs[i] = &TestObject{
				Int:   tag,
				Str:   strings.Repeat(strconv.Itoa(tag), 10),
				Float: float64(tag),
			}
		}
		return rs
	}
	createCollection := func(c context.Context, db *Collection, size int) {
		batchSize := 1000
		for i, length := 0, size/batchSize; i < length; i++ {
			cursor := i * batchSize
			objs := newObjects(cursor, cursor+batchSize)
			batch := db.coll.Bulk()
			for _, obj := range objs {
				batch.InsertOne(obj)
			}
			rs, err := batch.Run(c)
			if err != nil {
				b.Fatalf("create collection failed: %v", err)
			}
			if rs.InsertedCount != int64(len(objs)) {
				b.Fatalf("create collection failed: inserted:%d != %d", rs.InsertedCount, len(objs))
			}
		}
		count, err := db.coll.Find(context.TODO(), bson.M{}).Count()
		if err != nil {
			b.Fatalf("create collection failed: %v", err)
		}
		if count != int64(size) {
			b.Fatalf("create collection failed: documents.count(%d) != %d", count, size)
		}
	}
	for _, rowCount := range []int{1000, 10000, 10 * 10000, 100 * 10000} {
		database := "kv_benchmark"
		collection := fmt.Sprintf("TestObject_%d", rowCount)
		db, err := newMongoTest(database, collection)
		if err != nil {
			b.Fatalf("db error: %v", err)
		}
		b.Cleanup(func() {
			db.coll.DropCollection(context.TODO())
		})
		createCollection(context.TODO(), db, rowCount)

		name := fmt.Sprintf("Size=%d/", rowCount)
		b.ResetTimer()
		b.Run(name+"Set", func(b *testing.B) {
			obj := newObject(999)
			obj.Str = "1111111111"
			for i := 0; i < b.N; i++ {
				if err := db.Set(obj.Str, obj); err != nil {
					b.Fatalf("set failed: %v", err)
				}
			}
		})
		b.Run(name+"Get", func(b *testing.B) {
			obj := newObject(999)
			obj.Str = "111111111111"
			if err := db.Set(obj.Str, obj); err != nil {
				b.Fatalf("set failed: %v", err)
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if err := db.Get(obj.Str, obj); err != nil {
					b.Fatalf("get failed: %v", err)
				}
			}
		})
	}
}
