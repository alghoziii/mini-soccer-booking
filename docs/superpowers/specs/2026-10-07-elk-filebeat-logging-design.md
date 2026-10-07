# Design Spec: Centralized Logging with ELK Stack and Filebeat

## 1. Context and Goals
The project `mini-soccer-booking` consists of multiple microservices (`user-service`, `order-service`, `payment-service`, `field-service`) written in Go. Currently, logging is decentralized, with each service writing plain text logs to the console using `logrus`.

**Goal**: Implement a centralized logging architecture using Elasticsearch, Filebeat, and Kibana. As a proof of concept, this will be implemented for `user-service` first. 

## 2. Selected Approach
**Option B: Docker Logging Driver + Filebeat Autodiscover (Cloud-Native)**

Instead of having the Go application write to a physical file on a mounted volume, `user-service` will format its logs as JSON and output them to `stdout`. Filebeat will be configured with Docker Autodiscover to automatically detect the `user-service` container, read its Docker-managed console logs, and forward them directly to Elasticsearch. 

This approach offloads file management and log rotation to Docker, keeping the Go application lightweight.

## 3. Architecture & Data Flow
1. **user-service (Go)**: Logs events using `logrus` with `JSONFormatter` to `stdout`.
2. **Docker Engine**: Captures `stdout` and writes it to internal container log files (JSON format).
3. **Filebeat**: Uses autodiscover to mount and read the Docker container log files, decodes the JSON payload, and forwards it.
4. **Elasticsearch**: Stores and indexes the logs.
5. **Kibana**: Provides the UI to search and visualize the logs.

## 4. Components & Implementation Details

### A. Code Changes in `user-service-main`
- **Logger Initialization**: Instead of a new package, directly configure `logrus` at the very beginning of the `main()` function in `cmd/main.go` (or `main.go`):
  - Set formatter to `&logrus.JSONFormatter{}`.
  - Set output to `os.Stdout`.
  - Add a Hook to inject a default `"service": "user-service"` field.

### B. ELK Stack Infrastructure (`logging-management/docker-compose.yml`)
Create a new directory (e.g., `logging-management`) to host the logging infrastructure.
- **elasticsearch**: Single node, security disabled (for PoC). Port 9200.
- **kibana**: Exposes port 5601. Depends on elasticsearch.
- **filebeat**: Mounts `/var/lib/docker/containers` and `/var/run/docker.sock` to enable Docker Autodiscover. Mounts `filebeat.yml`.

### C. Filebeat Configuration (`filebeat.yml`)
```yaml
filebeat.autodiscover:
  providers:
    - type: docker
      hints.enabled: true
      templates:
        - condition:
            contains:
              docker.container.name: "user-service"
          config:
            - type: container
              paths:
                - /var/lib/docker/containers/${data.docker.container.id}/*.log
              # Decode the JSON payload from the Go app which is inside the Docker log wrapper
              processors:
                - decode_json_fields:
                    fields: ["message"]
                    target: "app_log"
                    overwrite_keys: true

output.elasticsearch:
  hosts: ["elasticsearch:9200"]
```


## 5. Testing Strategy
1. Start the ELK stack using `docker-compose up -d` in the `logging-management` directory.
2. Start the `user-service` using its `docker-compose.yaml` or `Dockerfile`.
3. Hit a `user-service` endpoint to generate a log.
4. Verify in Kibana (`http://localhost:5601`) that the log appears, is parsed as JSON, and contains the `"service": "user-service"` field.
