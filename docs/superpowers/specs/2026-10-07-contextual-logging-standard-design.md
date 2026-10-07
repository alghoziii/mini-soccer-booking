# Contextual Logging Standard Design Spec

## 1. Context & Motivation
Currently, `user-service` outputs JSON logs, but lacks business logging and request traceability. When an error occurs, it is difficult to trace which HTTP request caused it. We need a standardized logging approach that ties every log entry to a specific HTTP request using a `request_id`.

## 2. Goals
- Automatically generate a unique `request_id` for every incoming HTTP request.
- Log every HTTP request automatically (method, path, status, latency, request_id).
- Enable developers to use `logrus.WithContext(ctx)` to automatically inject the `request_id` into their log output without a custom wrapper package.

## 3. Architecture & Components

### A. Logrus Context Hook (`cmd/main.go`)
Instead of a custom wrapper package, we leverage `logrus.WithContext(ctx)`. We will implement a `logrus.Hook` in `cmd/main.go` that intercepts all logs, checks the Go standard `context.Context` for a `request_id`, and injects it into the log payload.

```go
type RequestIDHook struct{}

func (h *RequestIDHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *RequestIDHook) Fire(e *logrus.Entry) error {
	if e.Context != nil {
		if reqID, ok := e.Context.Value("request_id").(string); ok {
			e.Data["request_id"] = reqID
		}
	}
	return nil
}
```
*We will register this hook in the `Run()` function along with `logrus.SetFormatter`.*

### B. Gin Logging Middleware (`middlewares/logger.go`)
A middleware that intercepts requests, generates a UUID, stores it in the request's context, and logs the final response metrics.

```go
package middlewares

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"time"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Generate Request ID
		reqID := uuid.New().String()
		
		// Set in standard request context
		ctx := context.WithValue(c.Request.Context(), "request_id", reqID)
		c.Request = c.Request.WithContext(ctx)

		// Process request
		c.Next()

		// Log request details using the context
		latency := time.Since(start)
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"method":  c.Request.Method,
			"path":    c.Request.URL.Path,
			"status":  c.Writer.Status(),
			"latency": latency.String(),
		}).Info("HTTP Request")
	}
}
```

### C. Wiring it Up (`cmd/main.go`)
- Apply the `RequestLogger()` middleware to the Gin `router` in `cmd/main.go` (before the route groups).

## 4. Expected Outcome
Any developer can now do:
```go
// Inside user_service.go
logrus.WithContext(ctx).Info("Processing user login")
```
And the resulting Kibana log will automatically contain:
`{"level":"info","msg":"Processing user login","request_id":"123e4567-e89b-12d3-a456-426614174000","time":"..."}`
