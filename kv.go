package kv

import "github.com/cupen/kv/utils"

type DB interface {
	Get(key, val interface{}) error
	Set(key, val interface{}) error
	Del(key interface{}) error
	// SetIfExists(key string, val interface{}) error
	// SetIfNotExists(key string, val interface{}) error
}

func IsNotFound(err error) bool {
	return err == utils.ErrNotFound
}
