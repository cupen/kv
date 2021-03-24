package kv

import (
	"context"
	"fmt"
	urllib "net/url"

	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Config ...
type Config struct {
	KeySpace string
	Qmgo     *qmgo.Config
}

// MongoCollection ...
type MongoCollection struct {
	client *qmgo.Client
	db     *qmgo.Database
	coll   *qmgo.Collection
}

// NewMongoCollection ...
func NewMongoCollection(client *qmgo.Client, dbName, collName string) (*MongoCollection, error) {
	db := client.Database(dbName)
	return &MongoCollection{
		client: client,
		db:     db,
		coll:   db.Collection(collName),
	}, nil
}

// ConnectMongo ...
func ConnectMongo(conf *qmgo.Config) (*qmgo.Client, error) {
	var dbName = conf.Database
	if dbName == "" {
		u, err := urllib.Parse(conf.Uri)
		if err != nil {
			return nil, err
		}
		dbName = u.EscapedPath()[1:]
	}
	if dbName == "" {
		return nil, fmt.Errorf("Empty mongodb database name")
	}

	ctx := context.Background()
	client, err := qmgo.NewClient(ctx, conf)
	if err != nil {
		return nil, err
	}
	return client, err
}

// Get ...
func (mc *MongoCollection) Get(key, val interface{}) error {
	c := context.Background()
	q := mc.coll.Find(c, bson.M{"_id": key}).Limit(1)
	if err := q.One(val); err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrNotFound
		}
		return err
	}
	return nil
}

// Set ...
func (mc *MongoCollection) Set(key, val interface{}) error {
	c := context.Background()
	_, err := mc.coll.UpsertId(c, key, val)
	return err
}

// Del ...
func (mc *MongoCollection) Del(key interface{}) error {
	c := context.Background()
	if err := mc.coll.RemoveId(c, key); err != nil {
		if err == qmgo.ErrNoSuchDocuments {
			return nil
		}
		return err
	}
	return nil
}
