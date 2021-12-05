package edb

import "github.com/gotomicro/ego/core/elog"

// Container ...
type Container struct {
	config    *config
	name      string
	logger    *elog.Component
}
