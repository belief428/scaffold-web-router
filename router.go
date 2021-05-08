package WebRouter

import (
	"time"

	"github.com/belief428/scaffold-web-router/rate"

	"github.com/gin-gonic/gin"
)

type Router struct {
	Mode string
	*RouterConfig
}

type (
	RouterConfig struct {
		IPRate bool `json:"ip_rate"`
		*RouterIPLimitConfig
	}

	RouterIPLimitConfig struct {
		Limit    int `json:"limit"`
		Capacity int `json:"capacity"`
	}
)

type RouterServer func(config *RouterConfig) *Router

func (this *Router) Engine() *gin.Engine {
	// 版本
	gin.SetMode(this.Mode)
	// 创建gin日志
	app := gin.New()
	app.Use(Cors())
	app.NoRoute(NoRouteHandler())
	app.NoMethod(NoMethodHandler())
	app.Use(TimeoutHandle(time.Second * 1))
	app.Use(RecoveryHandler())

	if this.IPRate {
		if this.RouterIPLimitConfig == nil {
			panic("IP Limit Config Is Nil")
		}
		rate.NewIPRateLimiter()(this.Limit, this.Capacity)
		app.Use(rate.RequestIPRateLimiter()(this.Limit, this.Capacity))
	}
	return app
}

func NewRouter() RouterServer {
	return func(config *RouterConfig) *Router {
		return &Router{RouterConfig: config}
	}
}
