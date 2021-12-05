package eqmgo


import (
	"context"
	"github.com/qiniu/qmgo"
	opts "github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
)

type Database struct {
	mu        sync.Mutex
	db        *qmgo.Database
	processor processor
	logMode   bool
}

// Collection gets collection from database
func (d *Database) Collection(name string) *Collection {
	var cp *qmgo.Collection
	cp = d.db.Collection(name)

	return &Collection{
		coll: cp,
		processor: d.processor,
		logMode: d.logMode,
	}
}

// GetDatabaseName returns the name of database
func (d *Database) GetDatabaseName() string {
	return d.db.GetDatabaseName()
}

// DropDatabase drops database
func (d *Database) DropDatabase(ctx context.Context) error {
	return d.db.DropDatabase(ctx)
}

// RunCommand executes the given command against the database.
//
// The runCommand parameter must be a document for the command to be executed. It cannot be nil.
// This must be an order-preserving type such as bson.D. Map types such as bson.M are not valid.
// If the command document contains a session ID or any transaction-specific fields, the behavior is undefined.
//
// The opts parameter can be used to specify options for this operation (see the options.RunCmdOptions documentation).
func (d *Database) RunCommand(ctx context.Context, runCommand interface{}, opts ...opts.RunCommandOptions) *mongo.SingleResult {
	return d.db.RunCommand(ctx, runCommand, opts...)
}
