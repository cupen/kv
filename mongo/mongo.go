package mongo

import (
	"context"
	"fmt"
	"net/url"

	"github.com/cupen/kv/utils"
	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Options ...
type Options struct {
	URL string
}

// Collection ...
type Collection struct {
	// opts   *Options
	client *qmgo.Client
	db     *qmgo.Database
	coll   *qmgo.Collection
}

// New ...
func New(client *qmgo.Client, database, collection string) (*Collection, error) {
	db := client.Database(database)
	return &Collection{
		client: client,
		db:     db,
		coll:   db.Collection(collection),
	}, nil
}

// Connect ...
func Connect(_url string, opts ...options.ClientOptions) (*qmgo.Client, error) {
	u, err := url.Parse(_url)
	if err != nil {
		return nil, err
	}
	database := u.EscapedPath()[1:]
	if database == "" {
		return nil, fmt.Errorf("empty mongodb database name")
	}
	conf := qmgo.Config{
		Uri: _url,
	}
	ctx := context.Background()
	client, err := qmgo.NewClient(ctx, &conf, opts...)
	if err != nil {
		return nil, err
	}
	return client, err
}

// Get ...
func (mc *Collection) Get(key, val interface{}) error {
	c := context.Background()
	q := mc.coll.Find(c, bson.M{"_id": key}).Limit(1)
	if err := q.One(val); err != nil {
		if err == mongo.ErrNoDocuments {
			return utils.ErrNotFound
		}
		return err
	}
	return nil
}

// Set ...
func (mc *Collection) Set(key, val interface{}) error {
	c := context.Background()
	_, err := mc.coll.UpsertId(c, key, val)
	return err
}

// Del ...
func (mc *Collection) Del(key interface{}) error {
	c := context.Background()
	if err := mc.coll.RemoveId(c, key); err != nil {
		if err == qmgo.ErrNoSuchDocuments {
			return nil
		}
		return err
	}
	return nil
}
