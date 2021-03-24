package kv

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v7"
)

type RedisConfig struct {
	KeySpace string
	TypeName string
	TTL      time.Duration
}

// Redis ...
type Redis struct {
	client  *redis.Client
	baseKey string
	opts    *RedisConfig
}

// NewRedis ...
func NewRedis(client *redis.Client, opts *RedisConfig) (*Redis, error) {
	if client == nil {
		return nil, fmt.Errorf("nil client")
	}
	if opts == nil {
		return nil, fmt.Errorf("nil options")
	}
	return &Redis{
		client:  client,
		baseKey: fmt.Sprintf("%s:%s:", opts.KeySpace, opts.TypeName),
		opts:    opts,
	}, nil
}

func (r *Redis) buildKey(id interface{}) string {
	idtext := ""
	switch _id := id.(type) {
	case int:
		idtext = strconv.FormatInt(int64(_id), 10)
	case int32:
		idtext = strconv.FormatInt(int64(_id), 10)
	case int64:
		idtext = strconv.FormatInt(_id, 10)
	case string:
		idtext = _id
	default:
		panic(fmt.Errorf("Invalid id(%v) type<%t>", id, id))
	}

	if idtext == "" {
		panic(fmt.Errorf("Empty id(%v) type<%t>", id, id))
	}
	return r.baseKey + idtext
}

// Get ...
func (r *Redis) Get(id, val interface{}) error {
	key := r.buildKey(id)
	data, err := r.client.Get(key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return ErrNotFound
		}
		return err
	}
	return json.Unmarshal(data, val)
}

// Set ...
func (r *Redis) Set(id, val interface{}) error {
	key := r.buildKey(id)
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}
	return r.client.Set(key, data, r.opts.TTL).Err()
}

// Del ...
func (r *Redis) Del(id interface{}) error {
	key := r.buildKey(id)
	return r.client.Del(key).Err()
}
