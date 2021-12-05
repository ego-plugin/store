package eqmgo

import (
	"context"
	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
)

type Client struct {
	cc        *qmgo.Client
	processor processor
	database  string
	logMode   bool
}

func NewClient(ctx context.Context, conf *qmgo.Config, opts ...options.ClientOptions) (*Client, error) {
	var client *qmgo.Client
	var err error
	client, err = qmgo.NewClient(ctx, conf, opts...)
	if err != nil {
		return nil, err
	}
	return &Client{cc: client, processor: defaultProcessor, database: conf.Database}, nil
}

func (ec *Client) setLogMode(logMode bool) {
	ec.logMode = logMode
}

func defaultProcessor(processFn processFn) error {
	return processFn(&cmd{req: make([]interface{}, 0, 1)})
}

func (ec *Client) wrapProcessor(wrapFn func(processFn) processFn) {
	ec.processor = func(fn processFn) error {
		return wrapFn(fn)(&cmd{req: make([]interface{}, 0, 1)})
	}
}

func (ec *Client) Database(name string) *Database {
	var db *qmgo.Database
	_ = ec.processor(func(c *cmd) error {
		db = ec.cc.Database(name)
		logCmd(ec.logMode, c, "Database", db, name)
		return nil
	})
	if db == nil {
		return nil
	}
	return &Database{db: db, processor: ec.processor, logMode: ec.logMode}
}

func (ec *Client) DefaultDatabase() *Database {
	var db *qmgo.Database
	_ = ec.processor(func(c *cmd) error {
		db = ec.cc.Database(ec.database)
		logCmd(ec.logMode, c, "DefaultDatabase", db, ec.database)
		return nil
	})
	if db == nil {
		return nil
	}
	return &Database{db: db, processor: ec.processor, logMode: ec.logMode}
}

func (ec *Client) Session(opt ...*options.SessionOptions) (s *Session, err error) {
	var sess *qmgo.Session
	_ = ec.processor(func(c *cmd) error {
		sess, err = ec.cc.Session(opt...)
		logCmd(ec.logMode, c, "Session", nil)
		return err
	})
	return &Session{Session: sess, processor: ec.processor, logMode: ec.logMode}, err
}

func (ec *Client) Close(ctx context.Context) error {
	return ec.processor(func(c *cmd) error {
		logCmd(ec.logMode, c, "Close", nil)
		return ec.cc.Close(ctx)
	})
}

func (ec *Client) Ping(timeout int64) error {
	return ec.processor(func(c *cmd) error {
		logCmd(ec.logMode, c, "Ping", timeout, nil)
		return ec.cc.Ping(timeout)
	})
}

func (ec *Client) DoTransaction(ctx context.Context, callback func(sessCtx context.Context) (interface{}, error), opts ...*options.TransactionOptions) (x interface{}, err error) {
	return ec.cc.DoTransaction(ctx, callback, opts...)
}

// ServerVersion get the version of mongoDB server, like 4.4.0
func (ec *Client) ServerVersion() string {
	return ec.cc.ServerVersion()
}
