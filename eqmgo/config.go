package eqmgo

import (
	"time"
)

type config struct {
	// DSN DSN地址
	DSN string `json:"dsn" toml:"dsn"`
	// Debug 是否开启debug模式
	Debug bool `json:"debug" toml:"debug"`
	// DefaultDatabase 默认数据库
	DefaultDatabase string `json:"defaultDatabase" toml:"defaultDatabase"`
	// ConnectTimeoutMS specifies a timeout that is used for creating connections to the server.
	//	If set to 0, no timeout will be used.
	//	The default is 30 seconds.
	ConnectTimeoutMS int64 `json:"connectTimeoutMS" toml:"connectTimeoutMS"`
	// MaxPoolSize specifies that maximum number of connections allowed in the driver's connection pool to each server.
	// If this is 0, it will be set to math.MaxInt64,
	// The default is 100.
	MaxPoolSize uint64 `json:"maxPoolSize" toml:"maxPoolSize"`
	// MinPoolSize specifies the minimum number of connections allowed in the driver's connection pool to each server. If
	// this is non-zero, each server's pool will be maintained in the background to ensure that the size does not fall below
	// the minimum. This can also be set through the "minPoolSize" URI option (e.g. "minPoolSize=100"). The default is 0.
	MinPoolSize uint64 `json:"minPoolSize" toml:"minPoolSize"`
	// SocketTimeoutMS specifies how long the driver will wait for a socket read or write to return before returning a
	// network error. If this is 0 meaning no timeout is used and socket operations can block indefinitely.
	// The default is 300,000 ms.
	SocketTimeoutMS int64 `json:"socketTimeoutMS" toml:"socketTimeoutMS"`
	// EnableMetricInterceptor 是否启用prometheus metric拦截器
	EnableMetricInterceptor bool `json:"enableMetricInterceptor" toml:"enableMetricInterceptor"`
	// EnableAccessInterceptorReq 是否启用access req拦截器，此配置只有在EnableAccessInterceptor=true时才会生效
	EnableAccessInterceptorReq bool `json:"enableAccessInterceptorReq" toml:"enableAccessInterceptorReq"`
	// EnableAccessInterceptorRes 是否启用access res拦截器，此配置只有在EnableAccessInterceptor=true时才会生效
	EnableAccessInterceptorRes bool `json:"enableAccessInterceptorRes" toml:"enableAccessInterceptorRes"`
	// EnableAccessInterceptor 是否启用access拦截器
	EnableAccessInterceptor bool `json:"enableAccessInterceptor" toml:"enableAccessInterceptor"`
	// SlowLogThreshold 慢日志门限值，超过该门限值的请求，将被记录到慢日志中
	SlowLogThreshold time.Duration
	// interceptors 拦截器
	interceptors []Interceptor
}

// DefaultConfig 返回默认配置
func DefaultConfig() *config {
	c := new(config)
	c.DSN = "mongodb://root:example@127.0.0.1:27017"
	c.ConnectTimeoutMS = 30
	c.MaxPoolSize = 100
	c.MinPoolSize = 0
	c.SocketTimeoutMS = 30000
	return c
}
