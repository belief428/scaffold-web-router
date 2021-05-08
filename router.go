package WebRouter

import (
	"time"

	"github.com/gin-gonic/gin"
)

type Router struct{}

type RouterServer func() *Router

func (this *Router) Engine() *gin.Engine {
	// 版本
	gin.SetMode(gin.ReleaseMode)
	// 创建gin日志
	app := gin.New()
	app.Use(Cors())
	app.NoRoute(NoRouteHandler())
	app.NoMethod(NoMethodHandler())
	app.Use(TimeoutHandle(time.Second * 1))
	app.Use(RecoveryHandler())
	return app
}

func NewRouter() RouterServer {
	return func() *Router {
		return &Router{}
	}
}
