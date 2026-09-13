package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func Logger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Read and restore request body
		var reqBodyBytes []byte
		if c.Request.Body != nil {
			reqBodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBodyBytes))
		}

		// Wrap response writer to capture response body
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		fields := []zap.Field{
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", status),
			zap.Duration("latency", latency),
			zap.String("ip", c.ClientIP()),
		}

		// Log full context for 4xx and 5xx errors
		if status >= 400 {
			// Sanitize request body
			sanitizedReqBody := sanitizeJSON(reqBodyBytes)

			fields = append(fields, zap.String("request_body", string(sanitizedReqBody)))
			fields = append(fields, zap.String("response_body", blw.body.String()))

			if len(c.Errors) > 0 {
				err := c.Errors.Last()
				fields = append(fields, zap.String("error", err.Err.Error()))
				// Gin's error metadata can contain layer info if we attach it,
				// but at least we log the actual error string here.
				if err.Meta != nil {
					fields = append(fields, zap.Any("error_meta", err.Meta))
				}
			}

			logger.Error("request failed", fields...)
		} else {
			logger.Info("request completed", fields...)
		}
	}
}

func sanitizeJSON(body []byte) []byte {
	if len(body) == 0 {
		return body
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		// Not JSON or invalid, return as is (could be dangerous if it contains raw passwords,
		// but assuming API takes JSON)
		return body
	}

	if _, ok := data["password"]; ok {
		data["password"] = "[FILTERED]"
	}

	sanitized, _ := json.Marshal(data)
	return sanitized
}
