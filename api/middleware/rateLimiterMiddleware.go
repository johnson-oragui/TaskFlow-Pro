package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RateLimiterMiddleware(rdb *redis.Client, maxRequests int, duration time.Duration, type_ string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		var clientIp string
		var key string

		if type_ == "global" {
			clientIp = c.ClientIP()
			key = fmt.Sprintf("rateL:%s", clientIp)
		} else {
			path := c.Request.URL.Path
			userID := c.GetString("currentUserId")
			if userID == "" {
				// fallback to IP or reject
				userID = c.ClientIP()
			}
			key = fmt.Sprintf("rateL:%s:%s", userID, path)
		}

		count, err := rdb.Incr(ctx, key).Result()

		if err != nil {
			c.Header("X-RateLimit-Limit", strconv.Itoa(maxRequests))
			c.Header("X-RateLimit-Remaining", strconv.Itoa(maxRequests-int(count)))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message":     "An Unexpected Error Occurred",
				"status_code": 500,
			})
			return
		}

		if count == 1 {
			// Set expiration when key is first created
			rdb.Expire(ctx, key, duration)
		}

		if count > int64(maxRequests) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"message":     "rate limit exceeded, try again later",
				"status_code": 429,
			})
			return
		}

		c.Next()
	}
}
