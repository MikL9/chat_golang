package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
)

// KeyRedis key to extract redis client from context
const KeyRedis = "redis_client"

func Redis(rdb *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(KeyRedis, rdb)
		ctx.Next()
	}
}