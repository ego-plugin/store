package eqmgo

import (
	"github.com/gotomicro/ego/core/elog"
)

const PackageName = "component.eqmgo"

// Component 组成部分 client (cmdable and config)
type Component struct {
	config *config
	client *Client
	logger *elog.Component
}

// Client returns emongo Client
func (c *Component) Client() *Client {
	return c.client
}