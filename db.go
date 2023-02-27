package kv

import "github.com/cupen/kv/errors"

type DB interface {
	Get(key, val interface{}) error
	Set(key, val interface{}) error
	Del(key interface{}) error
}

func IsNotFound(err error) bool {
	return err == errors.ErrNotFound
}

// type DB[K comparable, T any] interface {
// 	Get(key K, val T) error
// 	Set(key K, val T) error
// 	Del(key K) error
// 	// SetIfExists(key string, val interface{}) error
// 	// SetIfNotExists(key string, val interface{}) error
// }

// func IsNotFound(err error) bool {
// 	return err == errors.ErrNotFound
// }

// type KV[K comparable, T any] struct {
// 	dbs []DB[K, T]
// }

// func New[K comparable, T any](dbs ...DB[K, T]) *KV[K, T] {
// 	return &KV[K, T]{
// 		dbs: dbs,
// 	}
// }
