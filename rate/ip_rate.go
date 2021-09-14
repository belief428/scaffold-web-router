package rate

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// IPRateLimiter IP限流器
type IPRateLimiter struct {
	ips      map[string]*rate.Limiter // IP地址库
	limit    int
	capacity int
	mutex    *sync.RWMutex // 锁
}

type AllowHandle func(ctx *gin.Context) bool

type IPRateLimiterConfig func(limit, capacity int) gin.HandlerFunc

var IPRateLimiterHandle *IPRateLimiter

// Handle
func (this *IPRateLimiter) Handle() AllowHandle {
	return func(ctx *gin.Context) bool {
		ip := ctx.ClientIP()

		this.mutex.RLock()
		defer this.mutex.RUnlock()

		limiter, exists := this.ips[ip]

		if !exists {
			limiter = rate.NewLimiter(rate.Limit(this.limit), this.capacity)
			this.ips[ip] = limiter
		}
		return limiter.Allow()
	}
}

// NewIPRateLimiter 初始化限流器
func NewIPRateLimiter() IPRateLimiterConfig {
	return func(limit, capacity int) gin.HandlerFunc {
		IPRateLimiterHandle = &IPRateLimiter{
			ips:      make(map[string]*rate.Limiter, 0),
			limit:    limit,
			capacity: capacity,
			mutex:    new(sync.RWMutex),
		}
		return nil
	}
}

// RequestIPRateLimiter 请求限流
func RequestIPRateLimiter() IPRateLimiterConfig {
	return func(limit, capacity int) gin.HandlerFunc {
		return func(c *gin.Context) {
			if IPRateLimiterHandle == nil {
				NewIPRateLimiter()(limit, capacity)
			}
			if !IPRateLimiterHandle.Handle()(c) {
				c.JSON(http.StatusOK, gin.H{
					"code": http.StatusTooManyRequests,
					"msg":  "访问频率过快，请稍后访问！",
				})
				c.Abort()
				return
			}
			c.Next()
		}
	}
}
