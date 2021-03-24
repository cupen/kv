package kv

import "errors"

var (
	ErrNotFound = errors.New("NotFound")
)

type DB interface {
	Get(key, val interface{}) error
	Set(key, val interface{}) error
	Del(key interface{}) error
	// SetIfExists(key string, val interface{}) error
	// SetIfNotExists(key string, val interface{}) error
}
