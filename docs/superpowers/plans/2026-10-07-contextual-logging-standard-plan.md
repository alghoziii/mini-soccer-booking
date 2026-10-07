# Contextual Logging Standard Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Implement a contextual logging standard that injects a `request_id` into every log entry across the application, leveraging Gin middleware and native logrus context hooks.

**Architecture:** A standard Gin middleware creates a UUID for every request. A global logrus Hook intercepts every log entry and injects the `request_id` if it exists in the Go standard `context.Context`.

**Spec:** `docs/superpowers/specs/2026-10-07-contextual-logging-standard-design.md`

## Global Constraints
- Do not create a new `logger` package wrapper. Use native `logrus.WithContext(ctx)`.
- Use `context.Context` (Go standard context) to pass the `request_id`, not just Gin's context.

---

### Task 1: Create the RequestID Logrus Hook

**Files:**
- Modify: `user-service-main/cmd/main.go`

**Interfaces:**
- Produces: A globally registered `RequestIDHook` for `logrus`.

- [ ] **Step 1: Define `RequestIDHook` in `main.go`**
Add the `RequestIDHook` struct and its `Levels` and `Fire` methods as defined in the spec.
The `Fire` method should safely check if `e.Context != nil` and if so, extract `"request_id"` (string) and assign it to `e.Data["request_id"]`.

- [ ] **Step 2: Register the Hook**
In the `Run()` function, below `logrus.SetOutput(os.Stdout)`, add:
`logrus.AddHook(&RequestIDHook{})`

- [ ] **Step 3: Compile check**
Run: `cd user-service-main && go build -o app ./cmd/main.go`
Expected: Successful build.

- [ ] **Step 4: Commit**
```bash
git add user-service-main/cmd/main.go
git commit -m "feat(user-service): add global logrus hook for request_id extraction"
```

### Task 2: Implement and Apply Gin Logging Middleware

**Files:**
- Modify: `user-service-main/middlewares/middleware.go`
- Modify: `user-service-main/cmd/main.go`

**Interfaces:**
- Produces: `func RequestLogger() gin.HandlerFunc` inside `middlewares`.
- Consumes: The Gin router in `main.go`.

- [ ] **Step 1: Add `RequestLogger` middleware**
In `middlewares/middleware.go`, append the `RequestLogger` function. It must generate a UUID, set it to the `Request().Context()`, call `c.Next()`, and then use `logrus.WithContext(ctx).WithFields(...)` to log the HTTP request metrics as defined in the spec. Ensure `"github.com/google/uuid"` is imported.

- [ ] **Step 2: Apply middleware in `main.go`**
In `cmd/main.go`, add `router.Use(middlewares.RequestLogger())` immediately after `router := gin.Default()`.

- [ ] **Step 3: Compile check**
Run: `cd user-service-main && go build -o app ./cmd/main.go`
Expected: Successful build.

- [ ] **Step 4: Commit**
```bash
git add user-service-main/middlewares/middleware.go user-service-main/cmd/main.go
git commit -m "feat(user-service): implement and apply request logger middleware"
```
