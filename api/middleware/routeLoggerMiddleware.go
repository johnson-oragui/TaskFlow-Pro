package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func RouteLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// capture request body
		var requestBody map[string]interface{}

		if c.Request.Body != nil {
			// read req body
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			_ = json.Unmarshal(bodyBytes, &requestBody)
		}

		userAgent := c.GetHeader("User-Agent")
		userIp := c.ClientIP()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Log Request
		log.Printf(
			"%s INFO Request received | name=RouteMiddlewareLogger user_ip=%s user_agent=%s path=%s method=%s payload=%v",
			time.Now().Format("2006-01-02 14:0405,000"),
			userIp,
			userAgent,
			path,
			method,
			sanitizePayload(requestBody),
		)

		// continue processing request
		c.Next()

		// calculate res time
		duration := time.Since(start).Seconds()

		statusCode := c.Writer.Status()

		currentUser, _ := c.Get("current_user_id")
		if currentUser == nil {
			currentUser = "Not Authenticated Guest"
		}

		log.Printf(
			"%s INFO Request completed | name=RouteMiddlewareLogger user_ip=%s user_agent=%s current_user=%v path=%s method=%s status_code=%d process_time=%.2fs",
			time.Now().Format("2006-01-02 15:04:05,000"),
			userIp, userAgent, currentUser, path, method, statusCode, duration,
		)

	}
}

func sanitizePayload(data map[string]interface{}) map[string]interface{} {
	if data == nil {
		return nil
	}

	copy := make(map[string]interface{})

	for k, v := range data {
		if k == "password" || k == "token" || k == "confirmPassword" {
			copy[k] = "************"
		} else if k == "device_info" {
			// nested masking
			if device, ok := v.(map[string]interface{}); ok {
				deviceCopy := make(map[string]interface{})
				for dk, dv := range device {
					if dk == "device_id" {
						deviceCopy[dk] = "**************"
					} else {
						deviceCopy[dk] = dv
					}
				}
				copy[k] = deviceCopy
			} else {
				copy[k] = v
			}
		} else {
			copy[k] = v
		}
	}

	return copy
}
