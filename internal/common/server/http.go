package server

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// RunHTTPServer 按服务名从配置中读取监听地址，再交给 RunHTTPServerOnAddr 真正启动。
func RunHTTPServer(serviceName string, wrapper func(router *gin.Engine)) {
	addr := viper.Sub(serviceName).GetString("http-addr")
	if addr == "" {
		// TODO: Warning log
	}
	RunHTTPServerOnAddr(addr, wrapper)
}

// RunHTTPServerOnAddr 创建 gin.Engine，并把路由注册动作交给 wrapper 回调。
func RunHTTPServerOnAddr(addr string, wrapper func(router *gin.Engine)) {
	// gin.New 不带默认中间件，适合教学阶段观察“哪些能力是手动加进去的”。
	apiRouter := gin.New()
	// wrapper 负责把 OpenAPI 生成的 handler 绑定到具体路由。
	wrapper(apiRouter)
	// 这一行不会修改已有路由，只是创建了一个没有被接住的 RouterGroup。
	apiRouter.Group("/api")
	if err := apiRouter.Run(addr); err != nil {
		panic(err)
	}
}
