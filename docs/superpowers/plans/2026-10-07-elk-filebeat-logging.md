# ELK + Filebeat Centralized Logging Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement a centralized logging architecture for `user-service` using Elasticsearch, Filebeat, and Kibana (EFK stack).

**Architecture:** We use Docker Autodiscover for Filebeat to directly capture the JSON stdout logs from the `user-service` container, forwarding them to a local Elasticsearch node, viewable via Kibana.

**Tech Stack:** Go (logrus), Docker, Elasticsearch, Filebeat, Kibana.

**Spec:** `docs/superpowers/specs/2026-10-07-elk-filebeat-logging-design.md`

## Global Constraints
- `user-service` logs must be formatted as JSON to `stdout`.
- The `user-service` container logs must be read via Filebeat's docker autodiscover mechanism.
- Logstash is EXPLICITLY excluded based on the ponytail review. Filebeat connects directly to Elasticsearch.

## Review Focus
- **Unparsed JSON:** If the log wrapper from Docker isn't stripped, Kibana will show nested JSON strings instead of queryable fields. We configure Filebeat `decode_json_fields` to fix this.
- **Service Name Missing:** If `user-service` logs don't include `"service": "user-service"`, tracking across multiple services later will fail. The Go logger initialization hook handles this.

---

### Task 1: Set up EFK Infrastructure

**Files:**
- Create: `logging-management/docker-compose.yml`
- Create: `logging-management/filebeat.yml`

**Interfaces:**
- Produces: Elasticsearch listening on `:9200`, Kibana on `:5601`. Filebeat monitoring `/var/lib/docker/containers`.

- [ ] **Step 1: Create `logging-management/filebeat.yml`**
Configures Filebeat to autodiscover `user-service` containers, read their log files, decode the JSON payload in the `message` field to `app_log`, and output directly to Elasticsearch `elasticsearch:9200`. Use the exact config from the spec.

- [ ] **Step 2: Create `logging-management/docker-compose.yml`**
Defines three services:
1. `elasticsearch` (v8.10.2, single-node, xpack.security.enabled=false, port 9200).
2. `kibana` (v8.10.2, port 5601, depends on elasticsearch).
3. `filebeat` (v8.10.2, mounts docker.sock, container log dir, and filebeat.yml, depends on elasticsearch).

- [ ] **Step 3: Run the stack to verify it starts**
Run: `cd logging-management && docker-compose up -d && sleep 15 && curl -s http://localhost:9200 | grep cluster_name`
Expected: Output showing the elasticsearch cluster name (verifying ES is up).

- [ ] **Step 4: Commit**
```bash
git add logging-management/
git commit -m "feat: add EFK infrastructure via docker-compose"
```

### Task 2: Configure JSON Logger in `main.go`

**Files:**
- Modify: `user-service-main/cmd/main.go` (or wherever `main()` is)

**Interfaces:**
- Produces: JSON-formatted stdout logs containing a `"service": "user-service"` field.

- [ ] **Step 1: Configure `logrus` in `main.go`**
At the very beginning of the `main()` function, configure the global `logrus` instance:
1. `logrus.SetFormatter(&logrus.JSONFormatter{})`
2. `logrus.SetOutput(os.Stdout)`
3. Add a custom `logrus.Hook` that injects `"service": "user-service"` into every entry.

- [ ] **Step 2: Verify `user-service` compiles and runs**
Run: `cd user-service-main && go build -o app ./cmd/main.go`
Expected: Successful build.

- [ ] **Step 3: Verify end-to-end integration (Manual Check)**
1. Start `user-service` via its docker-compose.
2. Hit an endpoint to trigger a log.
3. Check Kibana (`http://localhost:5601`) to ensure the log is ingested, parsed as JSON, and has `"service": "user-service"`.

- [ ] **Step 4: Commit**
```bash
git add user-service-main/cmd/main.go
git commit -m "feat(user-service): configure centralized json logger"
```
