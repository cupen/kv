package main

import (
	"log"
	"time"

	"github.com/cupen/kv"
	"github.com/cupen/kv/db/mongo"
	"github.com/cupen/kv/db/redis"
)

type TestObject struct {
	ID string
}

func main() {
	mongoClient := kv.Must(mongo.Connect("mongodb://root:root@127.0.0.1:27017/kv?authSource=admin"))
	db1 := kv.Must(mongo.New(mongoClient, "kv-exmaple-basic", "hello"))

	opts := &redis.Options{"kv_test", "TestObject", 1 * time.Minute}
	db2 := kv.Must(redis.NewWithURL("redis://127.0.0.1:6379/15", opts))

	// store := kv.New(db2, db1)
	store := kv.NewStore[string, TestObject](db1, db2, func(s string) (TestObject, error) {
		return TestObject{ID: s}, nil
	})
	val := TestObject{ID: "hello"}
	err := store.Set("123", &val)
	if err != nil {
		panic(err)
	}

	val2, err := store.Get("123")
	if err != nil {
		panic(err)
	}
	log.Printf("%v", val2)
}
