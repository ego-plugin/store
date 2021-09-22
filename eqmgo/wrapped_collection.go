package eqmgo

import (
	"context"
	"github.com/qiniu/qmgo"
	opts "github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/mongo"
)

type processor func(fn processFn) error
type processFn func(*cmd) error

type cmd struct {
	name string
	req  []interface{}
	res  interface{}
}

func logCmd(logMode bool, c *cmd, name string, res interface{}, req ...interface{}) {
	// 只有开启log模式才会记录req、res
	if logMode {
		c.name = name
		c.req = append(c.req, req...)
		//switch res := res.(type) {
		//case *qmgo.SingleResult:
		//	val, _ := res.DecodeBytes()
		//	c.res = val
		//default:
		c.res = res
		//}
	}
}

type Collection struct {
	coll      *qmgo.Collection
	processor processor
	logMode   bool
}

// Find find by condition filter，return QueryI
func (ec *Collection) Find(ctx context.Context, filter interface{}, opts ...opts.FindOptions) (que qmgo.QueryI) {
	_ = ec.processor(func(c *cmd) error {
		que = ec.coll.Find(ctx, filter, opts...)
		logCmd(ec.logMode, c, "Find", filter, nil)
		return nil
	})
	return que
}

// InsertOne insert one document into the collection
// If InsertHook in opts is set, hook works on it, otherwise hook try the doc as hook
// Reference: https://docs.mongodb.com/manual/reference/command/insert/
func (ec *Collection) InsertOne(ctx context.Context, doc interface{}, opts ...opts.InsertOneOptions) (result *qmgo.InsertOneResult, err error) {
	err = ec.processor(func(c *cmd) error {
		result, err = ec.coll.InsertOne(ctx, doc, opts...)
		logCmd(ec.logMode, c, "InsertOne", doc, result)
		return err
	})
	return result, err
}

// InsertMany executes an insert command to insert multiple documents into the collection.
// If InsertHook in opts is set, hook works on it, otherwise hook try the doc as hook
// Reference: https://docs.mongodb.com/manual/reference/command/insert/
func (ec *Collection) InsertMany(ctx context.Context, docs interface{}, opts ...opts.InsertManyOptions) (result *qmgo.InsertManyResult, err error) {
	err = ec.processor(func(c *cmd) error {
		result, err = ec.coll.InsertMany(ctx, docs, opts...)
		logCmd(ec.logMode, c, "InsertMany", docs, result)
		return err
	})
	return result, err
}

// Upsert updates one documents if filter match, inserts one document if filter is not match, Error when the filter is invalid
// The replacement parameter must be a document that will be used to replace the selected document. It cannot be nil
// and cannot contain any update operators
// Reference: https://docs.mongodb.com/manual/reference/operator/update/
// If replacement has "_id" field and the document is existed, please initial it with existing id(even with Qmgo default field feature).
// Otherwise, "the (immutable) field '_id' altered" error happens.
func (ec *Collection) Upsert(ctx context.Context, filter interface{}, replacement interface{}, opts ...opts.UpsertOptions) (result *qmgo.UpdateResult, err error) {
	err = ec.processor(func(c *cmd) error {
		result, err = ec.coll.Upsert(ctx, filter, replacement, opts...)
		logCmd(ec.logMode, c, "Upsert", filter, replacement, result)
		return err
	})
	return result, err
}

// UpsertId updates one documents if id match, inserts one document if id is not match and the id will inject into the document
// The replacement parameter must be a document that will be used to replace the selected document. It cannot be nil
// and cannot contain any update operators
// Reference: https://docs.mongodb.com/manual/reference/operator/update/
func (ec *Collection) UpsertId(ctx context.Context, id interface{}, replacement interface{}, opts ...opts.UpsertOptions) (result *qmgo.UpdateResult, err error) {
	err = ec.processor(func(c *cmd) error {
		result, err = ec.coll.UpsertId(ctx, id, replacement, opts...)
		logCmd(ec.logMode, c, "UpsertId", id, replacement, result)
		return err
	})
	return result, err
}

// UpdateOne executes an update command to update at most one document in the collection.
// Reference: https://docs.mongodb.com/manual/reference/operator/update/
func (ec *Collection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...opts.UpdateOptions) (err error) {
	return ec.processor(func(c *cmd) error {
		err = ec.coll.UpdateOne(ctx, filter, update, opts...)
		logCmd(ec.logMode, c, "UpdateOne", filter, update)
		return err
	})
}

// UpdateId executes an update command to update at most one document in the collection.
// Reference: https://docs.mongodb.com/manual/reference/operator/update/
func (ec *Collection) UpdateId(ctx context.Context, id interface{}, update interface{}, opts ...opts.UpdateOptions) (err error) {
	return ec.processor(func(c *cmd) error {
		err = ec.coll.UpdateId(ctx, id, update, opts...)
		logCmd(ec.logMode, c, "UpdateId", id, update)
		return err
	})
}

// UpdateAll executes an update command to update documents in the collection.
// The matchedCount is 0 in UpdateResult if no document updated
// Reference: https://docs.mongodb.com/manual/reference/operator/update/
func (ec *Collection) UpdateAll(ctx context.Context, filter interface{}, update interface{}, opts ...opts.UpdateOptions) (result *qmgo.UpdateResult, err error) {
	err = ec.processor(func(c *cmd) error {
		result, err = ec.coll.UpdateAll(ctx, filter, opts)
		logCmd(ec.logMode, c, "UpdateAll", filter, update, result)
		return err
	})
	return result, err
}

// ReplaceOne executes an update command to update at most one document in the collection.
// If UpdateHook in opts is set, hook works on it, otherwise hook try the doc as hook
// Expect type of the doc is the define of user's document
func (ec *Collection) ReplaceOne(ctx context.Context, filter interface{}, doc interface{}, opts ...opts.ReplaceOptions) (err error) {
	return ec.processor(func(c *cmd) error {
		err = ec.coll.ReplaceOne(ctx, filter, doc, opts...)
		logCmd(ec.logMode, c, "ReplaceOne", filter, doc)
		return err
	})
}

// Remove executes a delete command to delete at most one document from the collection.
// if filter is bson.M{}，DeleteOne will delete one document in collection
// Reference: https://docs.mongodb.com/manual/reference/command/delete/
func (ec *Collection) Remove(ctx context.Context, filter interface{}, opts ...opts.RemoveOptions) (err error) {
	return ec.processor(func(c *cmd) error {
		err = ec.coll.Remove(ctx, filter, opts...)
		logCmd(ec.logMode, c, "Remove", filter, nil)
		return err
	})
}

// RemoveId executes a delete command to delete at most one document from the collection.
func (ec *Collection) RemoveId(ctx context.Context, id interface{}, opts ...opts.RemoveOptions) (err error) {
	return ec.processor(func(c *cmd) error {
		err = ec.coll.RemoveId(ctx, id, opts...)
		logCmd(ec.logMode, c, "RemoveId", id, nil)
		return err
	})
}

// RemoveAll executes a delete command to delete documents from the collection.
// If filter is bson.M{}，all ducuments in Collection will be deleted
// Reference: https://docs.mongodb.com/manual/reference/command/delete/
func (ec *Collection) RemoveAll(ctx context.Context, filter interface{}, opts ...opts.RemoveOptions) (result *qmgo.DeleteResult, err error) {
	err = ec.processor(func(c *cmd) error {
		result, err = ec.coll.RemoveAll(ctx, filter, opts...)
		logCmd(ec.logMode, c, "RemoveAll", filter, result)
		return err
	})
	return result, err
}

// Aggregate executes an aggregate command against the collection and returns a AggregateI to get resulting documents.
func (ec *Collection) Aggregate(ctx context.Context, pipeline interface{}, opts ...opts.AggregateOptions) (agg qmgo.AggregateI) {
	_ = ec.processor(func(c *cmd) error {
		agg = ec.coll.Aggregate(ctx, pipeline, opts...)
		logCmd(ec.logMode, c, "Aggregate", pipeline, nil)
		return nil
	})
	return agg
}

// EnsureIndexes Deprecated
// Recommend to use CreateIndexes / CreateOneIndex for more function)
// EnsureIndexes creates unique and non-unique indexes in collection
// the combination of indexes is different from CreateIndexes:
// if uniques/indexes is []string{"name"}, means create index "name"
// if uniques/indexes is []string{"name,-age","uid"} means create Compound indexes: name and -age, then create one index: uid
func (ec *Collection) EnsureIndexes(ctx context.Context, uniques []string, indexes []string) (err error) {
	return ec.processor(func(c *cmd) error {
		err = ec.coll.EnsureIndexes(ctx, uniques, indexes)
		logCmd(ec.logMode, c, "EnsureIndexes", uniques, indexes)
		return err
	})
}

// CreateIndexes creates multiple indexes in collection
// If the Key in opts.IndexModel is []string{"name"}, means create index: name
// If the Key in opts.IndexModel is []string{"name","-age"} means create Compound indexes: name and -age
func (ec *Collection) CreateIndexes(ctx context.Context, indexes []opts.IndexModel) (err error) {
	return ec.processor(func(c *cmd) error {
		err = ec.coll.CreateIndexes(ctx, indexes)
		logCmd(ec.logMode, c, "CreateIndexes", indexes, nil)
		return err
	})
}

// CreateOneIndex creates one index
// If the Key in opts.IndexModel is []string{"name"}, means create index name
// If the Key in opts.IndexModel is []string{"name","-age"} means create Compound index: name and -age
func (ec *Collection) CreateOneIndex(ctx context.Context, index opts.IndexModel) (err error) {
	return ec.processor(func(c *cmd) error {
		err = ec.coll.CreateOneIndex(ctx, index)
		logCmd(ec.logMode, c, "CreateOneIndex", index, nil)
		return err
	})
}

// DropAllIndexes drop all indexes on the collection except the index on the _id field
// if there is only _id field index on the collection, the function call will report an error
func (ec *Collection) DropAllIndexes(ctx context.Context) (err error) {
	return ec.processor(func(c *cmd) error {
		err = ec.coll.DropAllIndexes(ctx)
		logCmd(ec.logMode, c, "DropAllIndexes", nil, nil)
		return err
	})
}

// DropIndex drop indexes in collection, indexes that be dropped should be in line with inputting indexes
// The indexes is []string{"name"} means drop index: name
// The indexes is []string{"name","-age"} means drop Compound indexes: name and -age
func (ec *Collection) DropIndex(ctx context.Context, indexes []string) (err error) {
	return ec.processor(func(c *cmd) error {
		err = ec.coll.DropIndex(ctx, indexes)
		logCmd(ec.logMode, c, "DropIndex", indexes, nil)
		return err
	})
}

// DropCollection drops collection
// it's safe even collection is not exists
func (ec *Collection) DropCollection(ctx context.Context) (err error) {
	return ec.processor(func(c *cmd) error {
		err = ec.coll.DropCollection(ctx)
		logCmd(ec.logMode, c, "DropCollection", nil, nil)
		return err
	})
}

// CloneCollection creates a copy of the Collection
func (ec *Collection) CloneCollection() (collection *mongo.Collection, err error) {
	_ = ec.processor(func(c *cmd) error {
		collection, err = ec.coll.CloneCollection()
		logCmd(ec.logMode, c, "CloneCollection", nil, nil)
		return err
	})
	return collection, err
}

// GetCollectionName returns the name of collection
func (ec *Collection) GetCollectionName() (str string) {
	_ = ec.processor(func(c *cmd) error {
		str = ec.coll.GetCollectionName()
		logCmd(ec.logMode, c, "GetCollectionName", nil, str)
		return nil
	})
	return str
}

// Watch returns a change stream for all changes on the corresponding collection. See
// https://docs.mongodb.com/manual/changeStreams/ for more information about change streams.
func (ec *Collection) Watch(ctx context.Context, pipeline interface{}, opts ...*opts.ChangeStreamOptions) (changeStream *mongo.ChangeStream, err error) {
	_ = ec.processor(func(c *cmd) error {
		changeStream, err = ec.coll.Watch(ctx, pipeline, opts...)
		logCmd(ec.logMode, c, "Watch", pipeline, nil)
		return err
	})
	return changeStream, err
}
