package edb

import (
	"github.com/gotomicro/ego/core/util/xtime"
	"time"
)

const (
	dialectMysql    = "mysql"
	dialectPostgres = "postgres"
)

// config options
type config struct {
	Dialect                    string        // 选择数据库种类，默认mysql
	DSN                        string        // DSN地址: mysql://root:secret@tcp(127.0.0.1:3306)/mysql?charset=utf8mb4&collation=utf8mb4_general_ci&parseTime=True&loc=Local&timeout=1s&readTimeout=3s&writeTimeout=3s
	Debug                      bool          // 是否开启调试，默认不开启，开启后并加上export EGO_DEBUG=true，可以看到每次请求，配置名、地址、耗时、请求数据、响应数据
	RawDebug                   bool          // 是否开启原生调试开关，默认不开启
	MaxIdleConns               int           // 最大空闲连接数，默认10
	MaxOpenConns               int           // 最大活动连接数，默认100
	ConnMaxLifetime            time.Duration // 连接的最大存活时间，默认300s
	SlowLogThreshold           time.Duration // 慢日志阈值，默认500ms
	EnableMetricInterceptor    bool          // 是否开启监控，默认开启
	EnableTraceInterceptor     bool          // 是否开启链路追踪，默认开启
	EnableDetailSQL            bool          // 记录错误sql时,是否打印包含参数的完整sql语句，select * from aid = ?;
	EnableAccessInterceptor    bool          // 是否开启，记录请求数据
	EnableAccessInterceptorReq bool          // 是否开启记录请求参数
	EnableAccessInterceptorRes bool          // 是否开启记录响应参数
}

// DefaultConfig 返回默认配置
func DefaultConfig() *config {
	return &config{
		Dialect: dialectPostgres,
		DSN: "",
		Debug: true,
		MaxIdleConns: 10,
		MaxOpenConns: 100,
		ConnMaxLifetime: xtime.Duration("300s"),
		SlowLogThreshold: xtime.Duration("500ms"),
		EnableMetricInterceptor: true,
		EnableTraceInterceptor: true,
	}
}