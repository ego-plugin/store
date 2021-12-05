package eqmgo


import (
	"context"
	"github.com/qiniu/qmgo"
	opts "github.com/qiniu/qmgo/options"
)

type Session struct {
	*qmgo.Session
	processor processor
	logMode   bool
}

func (ws *Session) StartTransaction(ctx context.Context, cb func(sessCtx context.Context) (interface{}, error), opts ...*opts.TransactionOptions) (result interface{}, err error) {
	_ = ws.processor(func(c *cmd) error {
		result, err = ws.Session.StartTransaction(ctx,cb)
		logCmd(ws.logMode, c, "StartTransaction", nil)
		return err
	})
	return result, err
}

// EndSession will abort any existing transactions and close the session.
func (ws *Session) EndSession(ctx context.Context) {
	_ = ws.processor(func(c *cmd) error {
		ws.Session.EndSession(ctx)
		logCmd(ws.logMode, c,"EndSession", nil)
		return nil
	})
}

// AbortTransaction aborts the active transaction for this session. This method will return an error if there is no
// active transaction for this session or the transaction has been committed or aborted.
func (ws *Session) AbortTransaction(ctx context.Context) error {
	return ws.processor(func(c *cmd) error {
		logCmd(ws.logMode, c, "AbortTransaction", nil)
		return ws.Session.AbortTransaction(ctx)
	})
}
