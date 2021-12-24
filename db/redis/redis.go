package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/cupen/kv/errors"
	redis "github.com/go-redis/redis/v8"
)

type Options struct {
	Keyspace string
	TypeName string
	TTL      time.Duration
}

// Redis ...
type Redis struct {
	client  *redis.Client
	baseKey string
	opts    *Options
}

var ctx = context.Background()

// New ...
func New(client *redis.Client, opts *Options) (*Redis, error) {
	if client == nil {
		return nil, fmt.Errorf("nil client")
	}
	if opts == nil {
		return nil, fmt.Errorf("nil options")
	}
	return &Redis{
		client:  client,
		baseKey: fmt.Sprintf("%s:%s:", opts.Keyspace, opts.TypeName),
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
		panic(fmt.Errorf("invalid id(%v) type<%t>", id, id))
	}

	if idtext == "" {
		panic(fmt.Errorf("empty id(%v) type<%t>", id, id))
	}
	return r.baseKey + idtext
}

// Get ...
func (r *Redis) Get(id, val interface{}) error {
	key := r.buildKey(id)
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return errors.ErrNotFound
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
	return r.client.Set(ctx, key, data, r.opts.TTL).Err()
}

// Del ...
func (r *Redis) Del(id interface{}) error {
	key := r.buildKey(id)
	return r.client.Del(ctx, key).Err()
}
