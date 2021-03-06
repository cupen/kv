package kv

import (
	"fmt"
	"testing"

	"github.com/qiniu/qmgo"
	"github.com/stretchr/testify/require"
)

type TestObject struct {
	Int   int
	Str   string
	Float float64
	Bytes []byte
}

func newMongoTest(dbName, collName string) (*MongoCollection, error) {
	c, err := ConnectMongo(&qmgo.Config{
		Uri: "mongodb://root:root@127.0.0.1:30001/kv?authSource=admin",
	})

	db, err := NewMongoCollection(c, dbName, collName)
	var _ DB = db
	return db, err
}

func TestMongoCollection(t *testing.T) {
	assert := require.New(t)

	db, err := newMongoTest("kv_test", "TestObject")
	assert.NoError(err)
	t.Cleanup(func() {
		db.Del(1)
		db.Del(2)
	})

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

	db, err := newMongoTest("kv_test_1kdocs", "TestObject")
	assert.NoError(err)
	t.Cleanup(func() {
		for i := 0; i < 1000; i++ {
			db.Del(i)
		}
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
